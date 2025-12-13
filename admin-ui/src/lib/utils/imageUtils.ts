// ============================================================================
// Image Utilities for AI Diagnosis
// ============================================================================

/**
 * Convert a File or Blob to base64 string
 */
export async function fileToBase64(file: File | Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      const result = reader.result as string;
      // Remove data URL prefix (e.g., "data:image/jpeg;base64,")
      const base64 = result.split(',')[1];
      resolve(base64);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

/**
 * Check if a file is an image
 */
export function isImageFile(file: File | { type?: string; mime_type?: string; file_type?: string }): boolean {
  const type = file.type || (file as any).mime_type || (file as any).file_type || '';
  return type.startsWith('image/');
}

/**
 * Check if a file is a video
 */
export function isVideoFile(file: File | { type?: string; mime_type?: string; file_type?: string }): boolean {
  const type = file.type || (file as any).mime_type || (file as any).file_type || '';
  return type.startsWith('video/');
}

/**
 * Extract frames from video for AI analysis
 * Returns array of base64-encoded frames
 */
export async function extractVideoFrames(videoFile: File, frameCount: number = 3): Promise<string[]> {
  return new Promise((resolve, reject) => {
    const video = document.createElement('video');
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    
    if (!ctx) {
      reject(new Error('Could not get canvas context'));
      return;
    }

    video.preload = 'metadata';
    video.src = URL.createObjectURL(videoFile);
    
    video.onloadedmetadata = () => {
      const duration = video.duration;
      const frames: string[] = [];
      let currentFrame = 0;

      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;

      const captureFrame = () => {
        if (currentFrame >= frameCount) {
          URL.revokeObjectURL(video.src);
          resolve(frames);
          return;
        }

        // Seek to frame position
        const time = (duration / (frameCount + 1)) * (currentFrame + 1);
        video.currentTime = time;
      };

      video.onseeked = () => {
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
        const base64 = canvas.toDataURL('image/jpeg', 0.8).split(',')[1];
        frames.push(base64);
        currentFrame++;
        captureFrame();
      };

      video.onerror = () => {
        URL.revokeObjectURL(video.src);
        reject(new Error('Error loading video'));
      };

      captureFrame();
    };
  });
}

/**
 * Compress image before sending to AI
 */
export async function compressImage(file: File, maxWidth: number = 1920, quality: number = 0.8): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    
    reader.onload = (e) => {
      const img = new Image();
      
      img.onload = () => {
        const canvas = document.createElement('canvas');
        let width = img.width;
        let height = img.height;

        // Calculate new dimensions
        if (width > maxWidth) {
          height = (height * maxWidth) / width;
          width = maxWidth;
        }

        canvas.width = width;
        canvas.height = height;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
          reject(new Error('Could not get canvas context'));
          return;
        }

        ctx.drawImage(img, 0, 0, width, height);
        const base64 = canvas.toDataURL('image/jpeg', quality).split(',')[1];
        resolve(base64);
      };

      img.onerror = () => reject(new Error('Error loading image'));
      img.src = e.target?.result as string;
    };

    reader.onerror = () => reject(new Error('Error reading file'));
    reader.readAsDataURL(file);
  });
}
