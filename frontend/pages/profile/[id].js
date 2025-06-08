import Layout from "@/components/layout";
import { ProfileNavbar } from "@/components/navbar";
import { useState, useEffect, useRef } from "react";
import Cropper from "react-easy-crop";
import getCroppedImg from "@/utils/cropImage";
import Feed from "@/components/feed";
import Image from "next/image";

export async function getServerSideProps(context) {
  const { req, params } = context;
  const profileIdFromUrl = params.id;

  // 取得 jwt_token
  const jwtToken = req.cookies.jwt_token;

  // SSR 驗證：沒有 jwt_token 直接導回首頁
  if (!jwtToken) {
    return {
      redirect: {
        destination: "/",
        permanent: false,
      },
    };
  }

  // 取得目前登入者 userId（建議用 /auth/status 取得，避免用 user_id cookie）
  let loggedInUserId = null;
  try {
    const statusRes = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/auth/status`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          cookie: `jwt_token=${jwtToken}`,
        },
      }
    );
    if (statusRes.ok) {
      const statusData = await statusRes.json();
      loggedInUserId = statusData.userID;
    }
  } catch (e) {
    // ignore
  }

  // fetch user profile data
  const profileRes = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_PROFILE_API}${profileIdFromUrl}`,
    {
      method: "GET",
      headers: {
        cookie: `jwt_token=${jwtToken}`,
      },
    }
  );
  if (!profileRes.ok) {
    return {
      notFound: true,
    };
  }
  const profileData = await profileRes.json();
  if (!profileData || !profileData.user_id) {
    return {
      notFound: true,
    };
  }

  // fetch user's posts
  const postsRes = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_POSTS_API}/${profileIdFromUrl}`,
    {
      method: "GET",
      headers: {
        cookie: `jwt_token=${jwtToken}`,
      },
    }
  );
  if (!postsRes.ok) {
    return {
      notFound: true,
    };
  }
  const postsData = await postsRes.json();
  if (!postsData || !Array.isArray(postsData)) {
    return {
      notFound: true,
    };
  }

  // fetch user's followers and following counts
  try {
    const followersRes = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_USERS_API}${profileIdFromUrl}/followers`,
      {
        method: "GET",
        headers: {
          cookie: `jwt_token=${jwtToken}`,
        },
      }
    );
    if (!followersRes.ok) {
      throw new Error("Failed to fetch followers");
    }
    const followersData = await followersRes.json();
    profileData.followersCount = followersData.length;

    const followingRes = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_USERS_API}${profileIdFromUrl}/following`,
      {
        method: "GET",
        headers: {
          cookie: `jwt_token=${jwtToken}`,
        },
      }
    );
    if (!followingRes.ok) {
      throw new Error("Failed to fetch following");
    }
    const followingData = await followingRes.json();
    profileData.followingCount = followingData === null ? 0 : followingData.length;
  } catch (error) {
    console.error("Error fetching followers or following:", error);
    // 如果有錯誤，則設置默認值
    profileData.followersCount = 0;
    profileData.followingCount = 0;
  }

  // get current user's following status
  let isFollowing = false;
  if (loggedInUserId && loggedInUserId !== profileIdFromUrl) {
    try {
      const followingRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_USERS_API}${loggedInUserId}/following`,
        {
          method: "GET",
          headers: {
            cookie: `jwt_token=${jwtToken}`,
          },
        }
      );
      if (followingRes.ok) {
        const followingData = await followingRes.json();
        isFollowing = followingData.some(
          (user) => user.user_id === profileIdFromUrl
        );
      }
    } catch (error) {
      console.error("Error checking following status:", error);
      // 如果有錯誤，則默認為未關注
      isFollowing = false;
    }
  }

  return {
    props: {
      profileId: profileIdFromUrl,
      currentUserId: loggedInUserId,
      profileData: {
        user_id: profileData.user_id,
        username: profileData.username || "",
        avatar_url: profileData.avatar_url || null,
        bio: profileData.bio || "The user is too lazy to write a bio",
      },
      postsData: postsData || [],
      isFollowing: isFollowing,
    },
  };
}

export default function Profile({
  profileId,
  currentUserId,
  profileData,
  postsData,
  isFollowing,
}) {
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
  const [uploadingDots, setUploadingDots] = useState(0);
  const [isFollowed, setIsFollowed] = useState(isFollowing);

  useEffect(() => {
    let interval;
    if (step === "uploading") {
      interval = setInterval(() => {
        setUploadingDots((dots) => (dots + 1) % 4);
      }, 500);
    } else {
      setUploadingDots(0);
    }
    return () => clearInterval(interval);
  }, [step]);

  // ... (其他既有函式, 例如 handleEditOrDone, useEffect, 等等，保持不變) ...
  const handleEditOrDone = async () => {
    if (editOnClick) {
      try {
        const res = await fetch(
          process.env.NEXT_PUBLIC_API_BASE_URL +
            process.env.NEXT_PUBLIC_PROFILE_API +
            `${currentUserId}` +
            "/bio",
          {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ bio }),
          }
        );
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
    if (isOwnProfile) {
      setShowEditAvatar(true);
    }
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
      const uploadApiUrl = process.env.NEXT_PUBLIC_UPLOAD_AVATAR2S3_URL;
      const presignedUrlResponse = await fetch(
        `${uploadApiUrl}?fileName=avatar.png&fileType=image/png&userId=${profileId}`,
        { method: "GET" }
      );

      if (!presignedUrlResponse.ok) {
        throw new Error("無法從伺服器獲取上傳連結");
      }

      const uploadData = await presignedUrlResponse.json();
      console.log("從上傳 API 收到的資料:", uploadData);

      // 增強：驗證收到的資料是否完整
      if (!uploadData.presignedUrl || !uploadData.key) {
        throw new Error(
          "從 API 收到的資料格式不正確，缺少 presignedUrl 或 key"
        );
      }
      const { presignedUrl, key, finalFileUrl } = uploadData;
      console.log("獲取的 Permanent URL:", finalFileUrl);

      // 步驟 3: 使用獲取的 Presigned URL，將圖片 Blob 直接 PUT 到 S3
      console.log("步驟 3: 正在上傳圖片到 S3...");
      const uploadToS3Response = await fetch(presignedUrl, {
        method: "PUT",
        body: croppedImageBlob,
        headers: { "Content-Type": "image/png" },
      });

      if (!uploadToS3Response.ok) {
        throw new Error("上傳圖片到 S3 失敗");
      }
      console.log("圖片上傳成功！物件 Key:", key);

      setDisplayUrl(finalFileUrl);
      setStep("done");

      let res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_PROFILE_API}${currentUserId}/avatar`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ avatar_url: finalFileUrl }),
        }
      );
      if (!res.ok) {
        throw new Error("更新頭像失敗，請稍後再試");
      }
      console.log("頭像更新成功！");
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

  // follow/unfollow 按鈕的處理函式
  const handleFollowClick = async () => {
  if (isOwnProfile) return;
  try {
    let res;
    if (!isFollowed) {
      // Follow
      res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_USERS_API}${profileId}/follow`,
        {
          method: "POST",
          credentials: "include",
        }
      );
      if (!res.ok) throw new Error("Follow 操作失敗，請稍後再試");
      setIsFollowed((prev) => !prev);
    } else {
      // Unfollow
      res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_USERS_API}${profileId}/unfollow`,
        {
          method: "POST",
          credentials: "include",
        }
      );
      if (!res.ok) throw new Error("Unfollow 操作失敗，請稍後再試");
      setIsFollowed((prev) => !prev);
    }
  } catch (error) {
    alert(`${isFollowed ? "Unfollow" : "Follow"} 操作失敗：${error.message}`);
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
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="size-12 bg-black/50 rounded-full text-white"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                        />
                      </svg>
                    </span>
                  </button>
                </div>
                <div>
                  <h1 className="text-xl font-semibold text-black">
                    {username}
                  </h1>
                  <div className="flex items-center gap-2 text-gray-600">
                    <span className="text-sm text-gray-600">#{profileId}</span>
                    <span className="text-sm text-gray-600">
                      Joined on 2023-10-01
                    </span>
                  </div>
                </div>
                <div className="text-sm flex flex-row flex-1 gap-2 text-gray-600 justify-evenly">
                  <p>Follwers: 0</p>
                  <p>Follwing: 0</p>
                </div>
                {!isOwnProfile && (
                  <button
                    onClick={handleFollowClick}
                    className={`px-4 py-2 rounded ${
                      isFollowed
                        ? "bg-white text-black border border-gray-400 hover:bg-gray-100"
                        : "bg-black text-white hover:bg-gray-600"
                    }`}
                  >
                    <div className="flex items-center gap-2">
                      <p>{isFollowed ? "Unfollow" : "Follow"}</p>
                      {isFollowed ? (
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={1.5}
                          stroke="currentColor"
                          className="size-4"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M6 18L18 6M6 6l12 12"
                          />
                        </svg>
                      ) : (
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={1.5}
                          stroke="currentColor"
                          className="size-4"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M12 4.5v15m7.5-7.5h-15"
                          />
                        </svg>
                      )}
                    </div>
                  </button>
                )}
                {isOwnProfile && (
                  <button onClick={handleEditOrDone}>
                    <div className="bg-black text-white hover:bg-gray-600 px-4 py-2 rounded ">
                      <div className="flex items-center gap-2">
                        <p>{editOnClick ? "Done" : "Edit bio"}</p>
                        {editOnClick ? (
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={1.5}
                            stroke="currentColor"
                            className="size-4"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M4.5 12.75l6 6 9-13.5"
                            />
                          </svg>
                        ) : (
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={1.5}
                            stroke="currentColor"
                            className="size-4"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125"
                            />
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
                  onChange={(e) => setBio(e.target.value)}
                />
              </div>
            </div>
            <Feed feedData={postsData} />
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
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m9.75 9.75 4.5 4.5m0-4.5-4.5 4.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                />
              </svg>
            </button>
            {step === "done" ? (
              <div className="flex flex-row items-center justify-center">
                <h2 className="text-lg font-semibold ">上傳成功！</h2>
              </div>
            ) : (
              <div className="flex flex-row items-center justify-between">
                <h2 className="text-lg font-semibold ">是否更換頭像？</h2>
              </div>
            )}

            {step === "" && (
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
            )}

            {step === "choose" && (
              <div className="flex flex-row justify-center items-center gap-4 mt-8">
                <input
                  type="file"
                  accept="image/*"
                  ref={fileInputRef}
                  onChange={handleFileChange}
                  className="hidden"
                />
                <button
                  className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                  onClick={handleCloseEditAvatar}
                >
                  Cancel
                </button>
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
                  onChange={(e) => setZoom(Number(e.target.value))}
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
              <div className="text-center">
                Uploading{".".repeat(uploadingDots)}
              </div>
            )}
            {step === "done" && (
              <div className="flex flex-col text-center gap-2 items-center">
                {displayUrl && (
                  <Image
                    src={displayUrl}
                    alt="Preview"
                    width={128}
                    height={128}
                    className="w-32 h-32 rounded-full mt-6"
                  />
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </Layout>
  );
}
