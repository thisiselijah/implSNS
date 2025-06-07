import Layout from "@/components/layout";
import { ProfileNavbar } from "@/components/navbar";
import { useState, useEffect, useRef } from "react";
import Cropper from "react-easy-crop";
import getCroppedImg from "@/utils/cropImage";
import Feed from "@/components/feed";

export async function getServerSideProps(context) {
    const { req, params } = context;
    const profileIdFromUrl = params.id;
    let loggedInUserId = req.cookies.user_id || null;

    // fetch user profile data
    const res = await fetch(`http://localhost:8080/api/v1/pages/profile/${profileIdFromUrl}`, { method: "GET" });
    if (!res.ok) {
        console.error("Failed to fetch profile data:", res.statusText);
        return {
            notFound: true,
        };
    }
    const profileData = await res.json();
    if (!profileData || !profileData.user_id) {
        console.error("Profile data is missing userId:", profileData);
        return {
            notFound: true,
        };
    }
    let avatar_access_key = profileData.avatar_access_key || null;
    const bio = profileData.bio || "The user is too lazy to write a bio";
    let viewableUrl = null;

    if (avatar_access_key) {
        // fetch avatar image URL
        const url = "https://xay0zgmfxk.execute-api.us-east-1.amazonaws.com/getS3ViewUrl";
        const viewUrlResponse = await fetch(`${url}?key=${encodeURIComponent(avatar_access_key)}`);

        if (!viewUrlResponse.ok) {
            throw new Error('無法獲取讀取連結');
        }

        let res = await viewUrlResponse.json();
        viewableUrl = res.viewableUrl;
        
    }


    return {
        props: {
            profileId: profileIdFromUrl,
            currentUserId: loggedInUserId,
            profileData: {
                user_id: profileData.user_id,
                username: profileData.username || "",
                avatar_url: viewableUrl,
                bio: bio,
            },

        },
    };
}

export default function Profile({ profileId, currentUserId, profileData }) {
    const isOwnProfile = profileId === currentUserId && currentUserId !== null;
    const [editOnClick, setEditOnClick] = useState(false);
    const initialBio = profileData.bio;
    const username = profileData.username;
    const [bio, setBio] = useState(initialBio);
    const [showEditAvatar, setShowEditAvatar] = useState(false);
    const [step, setStep] = useState("");
    const [preview, setPreview] = useState(null);
    const fileInputRef = useRef();
    const [crop, setCrop] = useState({ x: 0, y: 0 });
    const [zoom, setZoom] = useState(1);
    const [croppedAreaPixels, setCroppedAreaPixels] = useState(null);
    const [displayUrl, setDisplayUrl] = useState(profileData.avatar_url);

    // ... (其他既有函式, 例如 handleEditOrDone, useEffect, 等等，保持不變) ...
    const handleEditOrDone = async () => {
        if (editOnClick) {
            try {
                const res = await fetch(process.env.NEXT_PUBLIC_API_BASE_URL+process.env.NEXT_PUBLIC_PROFILE_API+`${currentUserId}`+'/bio', {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ bio }),
                });
                if (!res.ok) {
                    alert("上傳 bio 失敗，請稍後再試");
                    setBio(initialBio);
                }
            } catch (err) {
                console.error(err);
            }
        }
        setEditOnClick((prev) => !prev);
    };

    
    const handleAvatarEditClick = () => {
        setShowEditAvatar(true);
    };

    const handleCloseEditAvatar = () => {
        setShowEditAvatar(false);
        setStep("");
        setPreview(null);
    };

    const onCropComplete = (croppedArea, croppedAreaPixels) => {
        setCroppedAreaPixels(croppedAreaPixels);
    };

    const cropAndUpload = async () => {
        if (!croppedAreaPixels || !preview) {
            alert("裁切區域或圖片預覽不存在，請重試");
            return;
        }

        setStep("uploading");
        setDisplayUrl(null);

        try {
            // 步驟 1: 根據使用者裁切的區域，產生圖片 Blob
            console.log("步驟 1: 正在裁切圖片...");
            const croppedImageBlob = await getCroppedImg(preview, croppedAreaPixels);
            console.log("裁切完成，Blob:", croppedImageBlob);

            // 步驟 2: 從環境變數讀取 API URL，並請求「上傳」用的 Presigned URL
            console.log("步驟 2: 正在請求上傳連結...");
            const uploadApiUrl = "https://ixfso6it05.execute-api.us-east-1.amazonaws.com/avatar2s3";
            const presignedUrlResponse = await fetch(
                `${uploadApiUrl}?fileName=avatar.png&fileType=image/png&userId=${profileId}`,
                { method: 'GET' }
            );

            if (!presignedUrlResponse.ok) {
                throw new Error('無法從伺服器獲取上傳連結');
            }

            const uploadData = await presignedUrlResponse.json();
            console.log("從上傳 API 收到的資料:", uploadData);

            // 增強：驗證收到的資料是否完整
            if (!uploadData.presignedUrl || !uploadData.key) {
                throw new Error('從 API 收到的資料格式不正確，缺少 presignedUrl 或 key');
            }
            const { presignedUrl, key, finalFileUrl } = uploadData;
            console.log("獲取的 Permanent URL:", finalFileUrl);

            // 步驟 3: 使用獲取的 Presigned URL，將圖片 Blob 直接 PUT 到 S3
            console.log("步驟 3: 正在上傳圖片到 S3...");
            const uploadToS3Response = await fetch(presignedUrl, {
                method: 'PUT',
                body: croppedImageBlob,
                headers: { 'Content-Type': 'image/png' },
            });

            if (!uploadToS3Response.ok) {
                throw new Error("上傳圖片到 S3 失敗");
            }
            console.log("圖片上傳成功！物件 Key:", key);

            // 步驟 4: 使用剛剛上傳的 key，請求「讀取」用的 Presigned URL
            console.log("步驟 4: 正在獲取臨時讀取連結...");
            const viewApiUrl = "https://xay0zgmfxk.execute-api.us-east-1.amazonaws.com/getS3ViewUrl";
            const viewUrlResponse = await fetch(`${viewApiUrl}?key=${encodeURIComponent(key)}`);

            if (!viewUrlResponse.ok) {
                throw new Error('無法獲取讀取連結');
            }

            const { viewableUrl } = await viewUrlResponse.json();
            console.log("獲取讀取連結成功:", viewableUrl);

            // 步驟 5: 更新 UI，顯示上傳成功的圖片
            setDisplayUrl(viewableUrl);
            setStep("done");

            let res = await fetch(`http://localhost:8080/api/v1/pages/profile/${currentUserId}/avatar`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ avatar_access_key: key }),
            });
            if (!res.ok) {
                throw new Error("更新頭像失敗，請稍後再試");
            }
            console.log("頭像更新成功！");
            alert("頭像上傳成功！");
            setShowEditAvatar(false); // 關閉編輯頭像視窗

        } catch (error) {
            console.error("上傳流程出錯:", error);
            alert(`上傳失敗：${error.message}`);
            setStep("crop"); // 讓使用者可以重試
        }
    };

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setPreview(URL.createObjectURL(file));
            setStep("crop");
        }
    };


    return (
        <Layout pageTitle={isOwnProfile ? "Profile" : `Profile-${profileId}`}>
            <ProfileNavbar />
            {/* ... JSX 結構保持不變 ... */}
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 px-4 min-h-screen">
                <main className="flex flex-col col-start-4 col-end-10 ">
                    <div className="bg-white p-0.5 rounded-lg shadow">
                        <div className="flex flex-col p-4 border-b border-gray-200">
                            <div className="flex items-center gap-4 p-4">
                                <div className="relative group">
                                    <button onClick={handleAvatarEditClick}>
                                        <img
                                            alt="Profile Avatar"
                                            // 這裡最終應該要顯示從資料庫來的頭像
                                            src={displayUrl}
                                            className="inline-block size-12 rounded-full ring-1 ring-[#B6B09F]"
                                        />
                                        <span className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-12 bg-black/50 rounded-full text-white">
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                                            </svg>
                                        </span>
                                    </button>
                                </div>
                                <div>
                                    <h1 className="text-xl font-semibold text-black">
                                        {username}
                                    </h1>
                                    <div className="flex items-center gap-2 text-gray-600">
                                        <span className="text-sm text-gray-600">
                                            #{profileId}
                                        </span>
                                        <span className="text-sm text-gray-600">
                                            Joined on 2023-10-01
                                        </span>
                                    </div>
                                </div>
                                <div className="text-sm flex flex-row flex-1 gap-2 text-gray-600 justify-evenly">
                                    <p>
                                        Follwers: 0
                                    </p>
                                    <p>
                                        Follwing: 0

                                    </p>
                                </div>
                                {isOwnProfile && (
                                    <button onClick={handleEditOrDone}>
                                        <div className="bg-black text-white hover:bg-gray-600 px-4 py-2 rounded ">
                                            <div className="flex items-center gap-2">
                                                <p>{editOnClick ? "Done" : "Edit"}</p>
                                                {editOnClick ? (
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-4">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                                                    </svg>
                                                ) : (
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-4">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125" />
                                                    </svg>
                                                )}
                                            </div>
                                        </div>
                                    </button>
                                )}
                            </div>
                            <div className="flex items-center gap-4 p-4">
                                <textarea
                                    className="flex-1 p-2 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F] resize-none"
                                    disabled={!editOnClick}
                                    rows={1}
                                    value={bio}
                                    onChange={e => setBio(e.target.value)}
                                />

                            </div>

                        </div>
                        <Feed />

                    </div>

                </main>
            </div>
            {showEditAvatar && (
                <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 shadow-lg max-w-md w-full relative">
                        {/* ... Modal JSX 保持不變 ... */}
                        <button
                            className="text-black hover:text-black absolute top-1 right-1"
                            onClick={handleCloseEditAvatar}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                                <path strokeLinecap="round" strokeLinejoin="round" d="m9.75 9.75 4.5 4.5m0-4.5-4.5 4.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                            </svg>
                        </button>
                        <div className="flex flex-row items-center justify-between">
                            <h2 className="text-lg font-semibold ">是否更換頭像？</h2>

                        </div>

                        <div className="flex py-4 gap-4 justify-end">
                            <button
                                className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                                onClick={handleCloseEditAvatar}
                            >
                                Cancel
                            </button>
                            <button
                                className="px-4 py-2 bg-black text-white rounded hover:bg-gray-600"
                                onClick={() => setStep("choose")}
                            >
                                Yes
                            </button>
                        </div>
                        {step === "choose" && (
                            <div className="flex flex-col items-center gap-4">
                                <input
                                    type="file"
                                    accept="image/*"
                                    ref={fileInputRef}
                                    onChange={handleFileChange}
                                    className="hidden"
                                />
                                <button
                                    className="px-4 py-2 bg-black text-white rounded"
                                    onClick={() => fileInputRef.current.click()}
                                >
                                    選擇圖片
                                </button>
                            </div>
                        )}
                        {step === "crop" && preview && (
                            <div className="flex flex-col items-center gap-4">
                                <div className="relative w-64 h-64 bg-gray-200">
                                    <Cropper
                                        image={preview}
                                        crop={crop}
                                        zoom={zoom}
                                        aspect={1}
                                        cropShape="round"
                                        showGrid={false}
                                        onCropChange={setCrop}
                                        onZoomChange={setZoom}
                                        onCropComplete={onCropComplete}
                                    />
                                </div>
                                <input
                                    type="range"
                                    min={1}
                                    max={3}
                                    step={0.01}
                                    value={zoom}
                                    onChange={e => setZoom(Number(e.target.value))}
                                    className="w-48"
                                />
                                <button
                                    className="px-4 py-2 bg-green-600 text-white rounded"
                                    onClick={cropAndUpload}
                                >
                                    Crop and Upload
                                </button>
                            </div>
                        )}
                        {step === "uploading" && (
                            <div className="text-center">Uploading...</div>
                        )}
                        {step === "done" && (
                            <div className="text-center text-green-600">
                                <p>上傳成功！</p>
                                {displayUrl && (
                                    <img src={displayUrl} alt="上傳預覽" className="w-32 h-32 rounded-full mx-auto my-4" />
                                )}
                                <button className="ml-4 underline" onClick={handleCloseEditAvatar}>關閉</button>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </Layout>
    );
}