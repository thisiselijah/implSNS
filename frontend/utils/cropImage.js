export default function getCroppedImg(imageSrc, pixelCrop, outputSize = 512) {
    return new Promise((resolve, reject) => {
        const image = new window.Image();
        image.crossOrigin = "anonymous";
        image.src = imageSrc;
        image.onload = () => {
            const canvas = document.createElement("canvas");
            canvas.width = outputSize;
            canvas.height = outputSize;
            const ctx = canvas.getContext("2d");

            ctx.drawImage(
                image,
                pixelCrop.x,
                pixelCrop.y,
                pixelCrop.width,
                pixelCrop.height,
                0,
                0,
                outputSize,
                outputSize
            );

            canvas.toBlob((blob) => {
                if (!blob) {
                    reject(new Error("Canvas is empty"));
                    return;
                }
                resolve(blob);
            }, "image/png");
        };
        image.onerror = reject;
    });
}