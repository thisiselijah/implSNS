import Image from "next/image";
import { useState, useEffect, useRef } from "react";

export default function CreatePost({ onClose, avatar_url }) {
  const avatarUrl = avatar_url || "/user.png"; // Default avatar if none provided

  // const [MarkdownOnClick, setMarkdownOnClick] = useState(false);
  const [imagePreviews, setImagePreviews] = useState([]); // Array of {id: string, file: File, url: string}
  const [postText, setPostText] = useState(""); // State for textarea

  // Ref to hold current previews for unmount cleanup
  const previewsRef = useRef(imagePreviews);
  // Update ref on every render so unmount cleanup has the latest list
  previewsRef.current = imagePreviews;

  useEffect(() => {
    // Cleanup on unmount: revoke all Object URLs
    return () => {
      previewsRef.current.forEach((p) => {
        URL.revokeObjectURL(p.url);
        // console.log("Unmount: Revoked " + p.url);
      });
    };
  }, []); // Empty dependency array ensures this runs only on mount and unmount

  // const handleMarkdownClick = () => {
  //   setMarkdownOnClick((prev) => !prev);
  // };

  const handleUploadImage = () => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = "image/*";
    input.multiple = true;
    input.onchange = (event) => {
      const files = event.target.files;
      if (files && files.length > 0) {
        const newImageObjects = Array.from(files).map((file) => {
          const id = `${file.name}-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
          return { id, file, url: URL.createObjectURL(file) };
        });
        setImagePreviews((prevPreviews) => [...prevPreviews, ...newImageObjects]);
      }
    };
    input.click();
  };

  const handleRemovePreviewImage = (idToRemove) => {
    const imageToRemove = imagePreviews.find((img) => img.id === idToRemove);
    if (imageToRemove) {
      URL.revokeObjectURL(imageToRemove.url);
    }
    setImagePreviews((prevPreviews) => prevPreviews.filter((img) => img.id !== idToRemove));
  };

  // NEW: Refactored handlePost function to implement the full upload flow
  const handlePost = async () => {
    if (!postText && imagePreviews.length === 0) {
      alert("Please write something or add an image.");
      return;
    }

    let uploadedImageKeys = [];

    try {
      // --- STAGE 1: Get Presigned URLs from your Lambda API ---
      if (imagePreviews.length > 0) {
        console.log("Stage 1: Getting presigned URLs...");
        const filesToUpload = imagePreviews.map(p => ({
          fileName: p.file.name,
          fileType: p.file.type,
        }));

        const presignedUrlResponse = await fetch(`${process.env.NEXT_PUBLIC_UPLOAD_IMAGES2S3_URL}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ files: filesToUpload }),
        });

        if (!presignedUrlResponse.ok) {
          throw new Error(`Failed to get upload URLs. Status: ${presignedUrlResponse.status}`);
        }
        
        const { uploadData } = await presignedUrlResponse.json();
        console.log("Received upload data:", uploadData);

        // --- STAGE 2: Upload files directly to S3 ---
        console.log("Stage 2: Uploading files to S3...");
        const uploadPromises = imagePreviews.map((preview, index) => {
          const { presignedUrl } = uploadData[index];
          return fetch(presignedUrl, {
            method: 'PUT',
            headers: {
              'Content-Type': preview.file.type,
            },
            body: preview.file,
          });
        });

        const uploadResults = await Promise.all(uploadPromises);
        
        if (uploadResults.some(res => !res.ok)) {
            throw new Error('One or more files failed to upload to S3.');
        }
        
        console.log("All files successfully uploaded to S3.");
        uploadedImageKeys = uploadData.map(d => d.s3Key);
      }

      // --- STAGE 3: Create the post by calling your Go backend ---
      console.log("Stage 3: Creating post in backend...");

      const s3BucketName = process.env.NEXT_PUBLIC_AWS_S3_BUCKET_NAME;
      const awsRegion = process.env.NEXT_PUBLIC_AWS_S3_REGION;

      // 1. Transform S3 keys into full URLs and place them in the 'media' array.
      const mediaPayload = uploadedImageKeys.map(key => ({
        type: 'image',
        url: `https://${s3BucketName}.s3.${awsRegion}.amazonaws.com/${key}`
      }));
      
      // 2. Assemble the final payload with the correct field names: 'content', 'media', 'tags'.
      const postPayload = {
        author_id: "default-author-id", // Replace with actual author ID
        content: postText,
        media: mediaPayload,
        tags: [], // Add tags if needed
        location: null, // Add location if needed
      };

      console.log(`Submitting to Go backend`, postPayload);

      const createPostResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_POSTS_API}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(postPayload),
        credentials: 'include', // Ensure cookies are sent
      });

      if (!createPostResponse.ok) {
        throw new Error('Failed to create post in the backend.');
      }
      
      // --- FINAL: Success and cleanup ---
      alert("Post submitted successfully! (Simulated backend call)");
      setPostText("");
      setImagePreviews([]);
      if (onClose) onClose();

    } catch (error) {
      console.error("An error occurred during the post creation process:", error);
      alert(`Error: ${error.message}`);
    }
  };

  return (
    <div className="flex flex-col bg-white rounded-lg self-start">
      <div className="flex flex-row items-center border-b border-gray-200 w-full">
        <button className="p-4 pl-6" onClick={onClose}>
          Cancel
        </button>
        <div className="flex-1 text-center font-semibold text-black">
          New Post
        </div>
        <div className="p-4 pl-6 pr-6">
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
              d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z"
            />
          </svg>
        </div>
      </div>

      <div className="flex flex-row items-start p-2 border-b border-gray-200">
        <div className="p-2 flex-shrink-0">
          <Image
            src={avatar_url}
            alt="User Avatar"
            width={36}
            height={36}
            className="object-cover rounded-full ring-1 ring-offset-1 ring-[#B6B09F]"
          />
        </div>

        <div className="relative flex flex-col flex-1 max-w-full">
          <textarea
            placeholder="Got something to share today?"
            className="p-2 w-full text-gray-600 focus:outline-none focus:border-none rounded-md resize-none "
            rows={3}
            value={postText}
            onChange={(e) => setPostText(e.target.value)}

          />

          {/* Multiple Image Previews Section */}
          {imagePreviews.length > 0 && (
            <div className="p-2 flex flex-row gap-2 overflow-x-auto w-116">
              {imagePreviews.map((img) => (
                <div
                  key={img.id}
                  className="relative w-36 h-24 flex-shrink-0 border border-gray-300 rounded-md overflow-hidden group"
                >
                  <Image
                    src={img.url}
                    alt={`Preview ${img.file.name}`}
                    layout="fill"
                    objectFit="cover"
                    className="rounded-md"
                  />
                  <button
                    onClick={() => handleRemovePreviewImage(img.id)}
                    className="absolute top-0.5 right-0.5 bg-black opacity-0 text-white rounded-full p-0.5 w-5 h-5 flex items-center justify-center text-xs leading-none group-hover:opacity-75 transition-opacity"
                    title="Remove image"
                    aria-label="Remove image"
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
                        d="M6 18 18 6M6 6l12 12"
                      />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
          {/* End Multiple Image Previews Section */}

          <span className="flex px-2 gap-2.5 text-gray-600">
            <button
              className="hover:bg-gray-100 p-1 rounded transition-colors"
              onClick={handleUploadImage}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-5 text-gray-600"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                />
              </svg>
            </button>
            <button className="hover:bg-gray-100 p-1 rounded transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-5 text-gray-600"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
                />
              </svg>
            </button>

            <button className="hover:bg-gray-100 p-1 rounded transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M15.182 15.182a4.5 4.5 0 0 1-6.364 0M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z"
                />
              </svg>
            </button>

            <button className="hover:bg-gray-100 p-1 rounded transition-colors">
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
                  d="M5.25 8.25h15m-16.5 7.5h15m-1.8-13.5-3.9 19.5m-2.1-19.5-3.9 19.5"
                />
              </svg>
            </button>

            {/* <button
              className="hover:bg-gray-100 p-1 rounded transition-colors"
              title="Markdown formatting"
              onClick={handleMarkdownClick}
            >
              {!MarkdownOnClick ? (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  className="size-5 text-gray-600 "
                  fill="none"
                >
                  <rect
                    x="3"
                    y="5"
                    width="18"
                    height="14"
                    rx="2"
                    fill="white"
                    stroke="currentColor"
                    strokeWidth="1.5"
                  />
                  <path
                    d="M6 16V9L9 13L12 9V16"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M15 13l2 2l2-2"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M17 11v4"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                  />
                </svg>
              ) : (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  className="size-5 text-[#B6B09F]"
                  fill="none"
                >
                  <rect
                    x="3"
                    y="5"
                    width="18"
                    height="14"
                    rx="2"
                    fill="white"
                    stroke="currentColor"
                    strokeWidth="1.5"
                  />
                  <path
                    d="M6 16V9L9 13L12 9V16"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M15 13l2 2l2-2"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M17 11v4"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                  />
                </svg>
              )}
            </button> */}
          </span>
        </div>
      </div>
      <div className="flex flex-row items-center justify-end p-2 px-4">
        <button onClick={handlePost} className="p-1 px-4 bg-[#000000] hover:bg-gray-600 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors duration-150">
          Post
        </button>
      </div>
    </div>
  );
}
