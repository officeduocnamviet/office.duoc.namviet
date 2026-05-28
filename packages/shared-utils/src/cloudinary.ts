/**
 * Tiện ích Cloudinary
 */

/**
 * Tối ưu URL Cloudinary
 * Tự động thêm f_auto,q_auto để giảm dung lượng ảnh
 */
export function optimize_cloudinary_url(url: string | null | undefined): string {
  if (!url) return '';
  
  // Kiểm tra xem đây có phải là link Cloudinary hợp lệ không
  if (!url.includes('res.cloudinary.com')) {
    return url;
  }

  // Nếu link đã có tham số tối ưu rồi thì không làm gì thêm
  if (url.includes('f_auto') || url.includes('q_auto')) {
    return url;
  }

  // Chèn tham số tối ưu hóa vào sau đoạn '/upload/'
  // Cấu trúc URL Cloudinary thường là: https://res.cloudinary.com/<cloud_name>/image/upload/<version>/<public_id>
  const uploadToken = '/upload/';
  const uploadIndex = url.indexOf(uploadToken);

  if (uploadIndex !== -1) {
    const urlPart1 = url.substring(0, uploadIndex + uploadToken.length);
    const urlPart2 = url.substring(uploadIndex + uploadToken.length);
    
    // Thêm f_auto,q_auto/ vào giữa
    return urlPart1 + "f_auto,q_auto/" + urlPart2;
  }

  return url;
}

/**
 * Hàm hỗ trợ upload ảnh (từ Frontend) thông qua API Golang Backend.
 * @param file File ảnh người dùng chọn
 * @param backendUploadUrl URL của API Upload (Mặc định: http://localhost:8080/api/v1/upload)
 * @returns {Promise<string>} Đường dẫn secure_url trả về từ Cloudinary
 */
export async function uploadToCloudinary(file: File, backendUploadUrl?: string): Promise<string> {
  const defaultUrl = process.env.NODE_ENV === 'production' 
    ? 'https://namviet-official-150879831872.asia-southeast1.run.app/api/v1/upload'
    : 'http://localhost:8080/api/v1/upload';
  const urlToUse = backendUploadUrl || defaultUrl;
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch(urlToUse, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    throw new Error('Upload failed');
  }

  const data = await response.json();
  return optimize_cloudinary_url(data.url);
}
