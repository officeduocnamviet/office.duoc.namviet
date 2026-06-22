# 🚀 Mô tả tính năng NAM VIỆT Admin

## QUY TẮC CHO AI:

1. TRƯỚC KHI CODE: Bắt buộc dùng File System MCP đọc file này để nắm "Bức tranh toàn cảnh" tránh thiết kế DATABASE đụng độ.
2. SAU KHI CODE: Bắt buộc dùng file System MCP đánh dấu `[x]` vào tính năng đã làm xong và ghi chú kỹ thuật vào dòng "Cập nhật".

---

Ví dụ như sau:

## 📦 MODULE 1: QUẢN LÝ SẢN PHẨM & KHO

- [x] 1.1. CRUD Danh sách sản phẩm cơ bản (Đã xong)
- [ ] 1.2. Cấu hình Đơn vị tính đa tầng (Hộp/Vỉ/Viên)
  - _Mô tả:_ Chú ý logic quy đổi giá bán.
  - _Cập nhật:_ (Core/Aura sẽ tự điền vào đây khi xong)
- [ ] 1.3. Tính năng AI phân loại thuốc theo Quy định của Bộ Y Tế (Dùng Gemini)

## 💰 MODULE 2: HÓA ĐƠN & TÀI CHÍNH

- [ ] 2.1. Extension đồng bộ HĐ VAT từ Tổng cục Thuế
  - _Mô tả:_ Nhận mẻ 50 hóa đơn/lần. Yêu cầu có Redis Queue và Idempotency chống đúp.
  - _Cập nhật:_ (Chờ thực thi)
- [ ] 2.2. ...

## CÁC NHIỆM VỤ CẦN XÂY DỰNG:

- website và Mobile App web admin quản trị (bao gồm quản lý toàn bộ hệ thống, với các module như: Nhân sự; bán hàng; tài chính; Khám bệnh; ...)
- website và Mobile App web Bán Lẻ (đối tượng khách lẻ, theo dõi sức khỏe cho 1 cá nhân hoặc 1 gia đình, xây dựng tập khách hàng trung thành)
- website và Mobile App web Bán Sỉ (Đối tượng khách sỉ, tập trung vào cung cấp các giá trị như: Sự tiện lợi về mua hàng; công nợ và chính sách bán hàng rõ ràng; giao hàng nhanh chóng; ... )

## CHI TIẾT CÁC TÍNH NĂNG

### A. Các giao diện khởi tạo thông tin Lõi ban đầu

1. Giao diện CRUD thông tin Công Ty (hoặc nhiều công ty)

2. Giao diện CRUD thông tin các Chi Nhánh, các chi nhánh có khai báo là thuộc công ty nào. Vị trí Chi nhánh để phục vụ cho việc Điểm Danh của Nhân Viên. Ngoài ra, mỗi kho nên có thêm 1 Kho để quản lý sản phẩm Vật lý là "Quà tặng". Ví dụ: với Kho Tổng Công ty Nam Việt, sẽ nhìn ra 'Kho Hàng Hóa' và 'Kho Quà Tặng' tại kho này. Nhằm mục đích sau này tạo phiếu Kiểm Kê + Các chương trình riêng cho từng Chi Nhánh / Kho hàng.

3. Giao diện CRUD Phân Quyền. Cần được phân quyền chi tiết đến các ô (trường) thông tin. Ví dụ: Cho phép nhân viên A cập nhật Sản phẩm. Nhưng không cho phép Nhân Viên A xem được ô 'Giá vốn' của sản phẩm. Các Quyền Cơ Bản như sau:
   - Quyền được NHÌN THẤY NÚT TRUY CẬP GIAO DIỆN được chỉ định (Ví dụ: cho phép xem giao diện 'Danh Sach Sản Phẩm' )
   - Quyền được NHÌN THẤY CÁC THÔNG TIN CHI TIẾT TRONG 1 GIAO DIỆN ĐƯỢC CHỈ ĐỊNH (Ví dụ: cho phép nhìn thấy cột 'Tồn kho' của giao diện 'Danh sách sản phẩm'; cho phép nhìn thấy dữ liệu tại ô 'Giá Vốn' trong giao diện 'Chi Tiết Sản Phẩm')
   - Quyền được TẠO và GỬI THÔNG BÁO. Ví dụ: Tạo/Gửi thông báo tới Khách Bán Lẻ (push notification); Tạo/Gửi Thông Báo tới Khách Bán Sỉ; Tạo/Gửi Thông báo Nội Bộ; ..

4. Giao diện CRUD tài khoản đăng nhập cho User.

5. Giao diện CRUD profile (thông tin cá nhân + ảnh hồ sơ + ...) chi tiết của mỗi Nhân viên. Khi truy cập vào Giao diện này cần xem được: thông tin cơ bản; hồ sơ giấy tờ; quá trình làm việc + thăng tiến; lịch sử bảng lương; khóa học mà User này đã tham gia và hoàn thành; .v.v..

6. Giao diện CRUD TẠO MẪU IN. Đây là giao diện để khai báo các mẫu in (Ví dụ: mẫu in bán lẻ cỡ K80; Mẫu in Bán Sỉ cỡ A4; Mẫu hợp đồng Lao Động với nhân viên Kinh Doanh; Mẫu Hợp Đồng Lao Động với Nhân Viên Kế Toán; ...) + Sử dụng các {Biến} có sẵn để in ra. Ví dụ: Biến {TenCongTy}; biến {DiaChiCongTy}; biến {MaSoThueCongTy}; biến {TenKhachHang}; biến {TenNhanVien} ...

### B. Các giao diện CRUD liên quan đến NHÂN SỰ

7. Giao diện CRUD LOẠI HỢP ĐỒNG trong khối Nhân Sự. Ví dụ: khai báo "Hợp Đồng Thử Việc"; "Hợp đồng 6 tháng"; "Hợp đồng 1 năm"; ... Đây là Giao diện để khai báo các thông tin cơ bản, phục vụ cho việc theo dõi quá trình làm việc của Nhân Viên và Lộ trình Phát triển của Nhân Viên 1 cách tự động. Ví dụ: "Hợp đồng Thử Việc" quy định: Lương Cơ Bản là 3 tr/tháng + KPIs + Thưởng. Trong hợp đồng thử việc, Khi 'Nhân Viên' hoàn thành các KPIs + vượt qua khóa đào tạo "Văn Hóa" thì tự động chuyển sang "Hợp Đồng 6 Tháng" (có lương cao hơn, nhiều quyền hơn).

8. Giao diện CRUD CHÍNH SÁCH LƯƠNG. Khai báo các loại lương cơ bản. Ví dụ: Khai báo "Lương HĐ Thử Việc": 3 tr/tháng; "Lương HĐ 6 tháng": 3,5tr/tháng

9. Giao diện KPIs: Khai báo các KPIs cho User. Ví dụ: Nhân Viên Kinh Doanh A: KPIs đạt 100 triệu - lương 2 triệu; Nhân Viên Nhân Sự B: KPIs tuyển được 1 bác sĩ: Lương 5tr; ....

10. Giao diện CRUD KHÓA HỌC + ĐÀO TẠO + KIỂM TRA SAU HỌC. Là nơi khai báo các khóa học dành cho mọi đối tượng: ví dụ: Khóa học Nội Bộ (đào tạo văn hóa; đào tạo về sản phẩm nội bộ; đào tạo kỹ năng sử dụng Word/Excel; ...); khóa học dành cho Khách Sỉ (Đào tạo Luật Dược; Đào tạo về Thuế + Tài Chính; ...). Với khóa học Nội bộ (giống như các khóa đào tạo) thì mỗi khóa học được gắn với 1 user + loại hợp đồng .. Cụ thể. Có thể tự động chấm điểm khi User tham gia các khóa đào tạo và làm bài kiểm tra.

11. Giao diện CRUD "Tạo Ca Làm Việc" để khai báo các ca làm việc, tại các chi nhánh, để user hàng tuần vào đăng ký lịch làm việc cụ thể. Kết hợp Chấm công + mức độ Hoàn Thành KPIs + Lương Cơ Bản (theo hợp đòng) + Các công việc hoàn thành ..v.v để tính lương tự động.

12. Giao diện CRUD "Giao việc": là giao diện tạo công việc và giao cụ thể cho 1 nhân sự. Theo dõi toàn bộ quá trình của Công việc (Tạo mới => Giao Cho Ai => Đang thực hiện/Hoàn thành => kết Quả công việc hiện tại => Quản lý Nghiệm Thu). Yêu cầu công việc có thể Đo và chấm điểm cụ thể, vì liên quan đến việc Tính lương.

### C. Các giao diện CRUD ban đầu liên quan đến TÀI CHÍNH

11. Giao diện CRUD "DANH SÁCH NGÂN HÀNG". đây là giao diện để sử dụng API của VietQR đồng bộ danh sách các Ngân Hàng đang hoạt động tại Việt Nam. Để phục vụ cho việc: Tạo Mã QR thanh toán; Quản lý thông tin Ngân Hàng toàn hệ thống một cách đồng bộ.

12. Giao diện CRUD "HỆ THỐNG TÀI KHOẢN KẾ TOÁN theo THÔNG TƯ 200". Để khai báo các số tài khoản Kế Toán, phục vụ cho việc Báo Cáo Tài Chính, Hạch toán, tạo bút toán ... tự động. Ví dụ: Mã tài khoản '111' tương ứng là 'Tiền mặt'; loại tài khoản 'Tài Sản'; tính chất 'Dư nợ'; trạng thái ...

13. Giao diện CRUD "Loại Phiếu Thu / Chi". Giao diện khai báo các loại phiếu thu/chi của công ty, được liên kết đến tài khoản Kế Toán được tạo ở trên.

14. Giao diện CRUD "QUẢN LÝ TÀI KHOẢN & QUỸ TIỀN" - Khai báo các quỹ tiền mặt và tài khoản ngân hàng của công ty, cần gắn với công ty nào (hoặc chi nhánh nào)

15. Giao diện CRUD "GIAO DỊCH LẶP LẠI" - Đây là giao diện tạo các phiếu Thu/Chi ĐỊNH KỲ và tự sinh ra. Ví dụ: Tạo 1 phiếu CHI, cứ đến ngày mùng 5 hàng tháng là thanh toán Bảo Hiểm Xã Hội qua ngân hàng VietcomBank, với thông tin: Hình thức Chuyển khoản; Người nhận: BHXH Việt Nam; Ngân hàng nhận: VietcomBank; Số TK: 111.222.333; Số tiền: 5.000.000đ; Lý do Chi: chi Bảo Hiểm Xã Hội

16. Giao diện CRUD "QUẢN LÝ TÀI SẢN" : khai báo các loại Tài Sản; tên Tài Sản; Tài Chính & Khấu Hao (để tính vào Lãi / Lỗ tại mỗi chi nhánh hoặc công ty sau này); Phân Bổ cho ai Sử dụng, ở Chi Nhánh/Công ty nào để tính Lãi/Lỗ cho đơn vị đó; Lịch bảo trì + Bảo Dưỡng ....

### D. Các giao diện CRUD ban đầu liên quan đến KẾT NỐI:

17. Giao diện CRUD "Bản Đồ Hành Chính Việt Nam". Nơi khai báo: Tỉnh; Huyện; Xã/Phường ....

18. Giao diện CRUD "Đơn Vị Vận Chuyển": Nơi Khai báo các đơn vị vận chuyển, ví dụ như: Xe Nội Bộ Công Ty (Lịch đi xe như nào); Viettel Post (Username và PassWord để kết nối); Shopee Express (có được cài đặt là Đơn Vị Vận Chuyển mặc định hay không); Xe Khách (gồm những xe nào; thời gian Vận Chuyển như nào; Cước Phí như nào ...)

19. Kết Nối Facebook: để kết nối và quản lý các Page. Ví dụ: Kết nối và quản lý Chat trên Page; tạo bài viết tự động; Lấy dữ liệu Quảng cáo về hệ thống (lượt Like; Lượt theo dõi mới; ...)

20. Kết Nối Zalo OA: tạo lịch bắn tin nhắn cho khách; trả lời khách trên Zalo OA

21. Kết nối Thương Mại Điện Tử: Ví dụ như sàn Shopee; ... Đồng bộ đơn hàng; trả lời Khách hàng; ...

22. Kết nối Tiktok (Nếu được)

23. Kết Nối đến Cục Dược => vì bán Thuốc cần liên thông cục Dược Quốc Gia để báo cáo. Hoặc tính năng xuất file Excel để báo cáo cục dược thông qua phần mềm offline báo cáo cục dược là phần mềm "Ánh Sáng" (khuyến khích giải pháp này vì Phần mềm này đã mua)

### E. Các giao diện CRUD ban đầu liên quan đến 'KẾT NỐI: AI & Công nghệ'

21. Cấu Hình Chatbot: khai báo ChatBot; kết nối đến API ChatBot; Chatbot áp dụng cho Khách Lẻ; ChatBot Áp dụng cho khách Sỉ; ChatBot nội bộ (để thực hiện các quyền thao tác với hệ thống nội bộ ..) => kéo toàn bộ lịch sử Chat về hệ thống để User Chăm Sóc Khách Hàng có thể theo dõi hoặc tiếp tục CSKH.

22. Cấu Hình Theme Website : Cho phép thay đổi Theme website theo chủ đề tùy từng Dịp lễ trong năm

23. Cấu hình Theme Mobile App : Cho phép thay đổi Theme Mobile App theo chủ đề tùy từng Dịp lễ trong năm

24. Quản lý Pop Up trên Website và Mobile App

25. Cấu Hình Chung khác (nếu cần thiết)

### F. Cấu Hình Chăm Sóc Khách Hàng:

26. Quy tắc 'Tích Điểm' : cách Tính điểm như nào; đổi điểm được quà gì; hạn sử dụng của Điểm; Áp dụng cho ai; Áp dụng tại Chi Nhánh và Công ty nào; ...
27. Quy Tắc 'Hạng Khách Hàng': giao diện Khai báo và ra luật thăng hạng cho khách hàng + các chương trình ưu đãi của mỗi hạng.

### G. Cấu Hình Nhà Cung Cấp

28. Giao diện CRUD "Nhà Cung Cấp": Quản lý các thông tin chung của hiện tại như: Tên; Mã Số Thuế; Giấy phép Kinh Doanh; Thời gian giao hàng của NCC; .v.v.

29. Giao diện CRUD "Đơn vị Bảo Hành" : phục vụ cho việc Tạo đơn bảo hành + Theo dõi quản lý

30. Giao diện "Ánh Xạ Sản Phẩm từ Nhà Cung Cấp" vào Sản phẩm trên hệ thống.

31. Giao diện "Hợp Đồng & CTKM": nơi khai báo các chương trình Khuyến Mại và Hợp Đồng của Nhà Cung Cấp. Sau này, mỗi khi mua hàng, chỉ việc bấm '1 nút' là hệ thống tự động tạo các đơn mua hàng cần mua, tự điền trước các chương trình KM của Nhà Cung Cấp để tính ra giá vốn nhập hàng cuối cùng + theo dõi các Chiết khấu trả sau của NCC.

### H. Cấu Hình Sản Phẩm:

32. Giao diện "Danh Sách Sản phẩm Vật Lý"
33. Giao diện "Chi tiết sản phẩm": quản lý các thông tin của sản phẩm như: Tên SP; Mã SKU; các loại đơn vị tính; Hình ảnh; ..v.v.
34. Giao diện "Sửa Quy Cách Sản phẩm Nhanh"
35. Giao diện "Sửa Giá Sản phẩm Nhanh": cập nhật lại giá vốn của sản phẩm cho từng lô; giá bán lẻ; giá bán buôn ...
36. Giao diện "Cài đặt Số lượng tồn kho Min / Max và Nhà Cung Cấp" của mỗi sản phẩm.
37. Giao diện "Cài đặt nhanh Vị Trí Sản Phẩm tại mỗi Kho"
38. Giao diện "Cập nhật Mã Vạch của sản phẩm"
39. Giao diện CRUD "Danh Mục Sản Phẩm": phân loại chi tiết sản phẩm, ví dụ như: Thuốc không Kê Đơn (OTC); Thuốc Kê Đơn (ETC); Thuốc Kê Đơn Thuộc Danh Mục Thuốc Thiết Yếu; Thuốc Kê Đơn thuộc Danh Mục Kiểm Soát Đặc Biệt; Thiết Bị Y Tế Loại A; Thiết Bị Y Tế loại B; Thiết Bị Y Tế Loại C; Thiết Bị Y Tế loại D; Thực Phẩm Chức Năng; Thực Phẩm Chức Năng Vitamin; Thực Phẩm Chức Năng Não; Thực Phẩm Chức Năng Tim Mạch; ..v.v.v (các nhóm Thực Phẩm Chức Năng còn lại); Dịch vụ Web Bán Hàng; ... Vắc Xin; Quà tặng (tức là không được bán, giá trị thanh toán là 0 đ)
40. Giao diện CRUD "Hoạt Chất Sản phẩm": Gắn 1 sản phẩm với 1 hoặc nhiều hoạt chất.
41. Giao diện Cài đặt tích chọn sản phẩm này sẽ được hiển thị bán trên các kênh nào: Áp Dụng Bán Sỉ (sẽ cho hiển thị trên Web/App bán sỉ); Áp dụng Bán Lẻ (sẽ cho hiển thị trên Web/App bán Lẻ); Bán trên Sàn TMĐT; Bán Trên Facebook; ....

### H. Cấu Hình Dịch Vụ và Bệnh Học:

41. Giao diện CRUD "Danh sách Bệnh" (theo bảng DIC10 của bộ y tế, nhưng cần rút gọn lại), đồng thời khai báo các triệu chứng điển hình của mỗi bệnh.
42. Giao diện CRUD "Đơn Thuốc Mẫu": gắn 1 hoặc nhiều Thuốc với 1 bệnh cụ thể. Có chia theo độ tuổi + hướng dẫn sử dụng cho các thuốc đó, phù hợp với từng bệnh và độ tuổi tương ứng.
43. Giao diện CRUD "Danh sách Mũi Tiêm Chủng": Tạo và quản lý danh sách các mũi Tiêm Chủng, mỗi mũi tiêm chủng gắn với 1 hoặc nhiều loại Vắc Xin.
44. Giao diện CRUD "Phác Đồ Tiêm Chủng": nơi tạo, quản lý và ra quy tắc cho 1 hoặc nhiều mũi Tiêm Chủng. Để Điều Dưỡng và Bác Sĩ không được thực hiện sai.

## SANG PHẦN LOGIC PHỤC VỤ CHO VIỆC HOẠT ĐỘNG và KINH DOANH:

45. Giao diện CRUD: Phân loại Khách Hàng và Quyền Mua/Bán hàng của khách hàng: Ví dụ: Khách hàng là cửa hàng Mẹ Bé sẽ không được nhìn thấy sản phẩm Thuốc trên web sỉ; Khách hàng là Quầy Thuốc không được bán thuốc thuộc danh mục "Thuốc Kê Đơn thuộc danh mục Thuốc Kiểm Soát Đặc Biệt"; ..v.v
46. Giao diện CRUD "Khách hàng Kênh Sỉ": nơi khai báo các khách hàng mua Sỉ (Đối tượng là các Nhà Thuốc; Quầy Thuốc; Phòng Khám; Cửa Hàng Mẹ Bé; ...). Họ cũng là các đơn vị Kinh Doanh. Đây là nơi hiển thị danh sách khách hàng sỉ và theo dõi các thông tin chung của khách hàng, ví dụ như: Ngày mua gần nhất; Nợ hiện tại; Thuốc phân loại Khách Hàng Sỉ nào; ...
47. Giao diện "Chi Tiết Khách Hàng Sỉ": là nơi cập nhật các thông tin chi tiết của 1 khách hàng.
48. Giao diện CRUD "Khách hàng Kênh Lẻ": nơi khai báo các khách hàng mua Lẻ (mua tại cửa hàng; mua tại website; Sử dụng dịch vụ khám bệnh; ...), theo dõi các thông tin chung của khách hàng, ví dụ như: Ngày mua gần nhất; Số điểm tích lũy; Mối quan hệ với Khách khác (Ví dụ: Khách A là mẹ của bé C; Khách D là ông nội của khách F); xem lịch sử giao diện; ..

## GIAO DIỆN MARKETING và CSKH:

49. Giao diện "Tạo nhóm Khách Hàng": là giao diện để tạo ra các nhóm khách hàng có cùng 1 đặc điểm. Với nhiều kịch bản khác nhau, ví dụ:

        - Nhóm khách hàng Bán Lẻ có cùng ngày sinh nhật tháng 6
        - Nhóm khách hàng Bán Sỉ giới tính Nữ - độ tuổi từ 20 đến 30 tuổi
        - Nhóm khách hàng Bán Lẻ có giới tính Nam, ngày mua cuối cùng cách đây 40 ngày
        - Nhóm khách hàng Bán Lẻ có cùng bệnh mãn tính 'Cao huyết Áp'
        - Nhóm khách hàng Bán Sỉ, Phân nhóm 'Mẹ Bé', nằm tại huyện ABC thuộc tỉnh XYZ

    Với việc tạo nhóm khách hàng này, có thể tạo ra các chương trình khuyến mại riêng. Tối ưu chi phí Tin nhắn Zalo hoặc Marketing cho từng nhóm
    Yêu cầu với giao diện này như sau:
    - Giao diện List phân khúc khách hàng: ngoài nút tạo lẻ từng phân khúc. Cần có nút 'Tạo Hàng Loạt' dựa vào các điều kiện khách hàng đã có sẵn, tự động tạo ra các nhóm khách hàng tương ứng.
    - Các nhóm khách hàng có khả tự động Thêm / Xóa bớt khách hàng trong nhóm. Ví dụ: nhóm khách hàng có sinh nhật tháng 7, hôm nay bán hàng có thêm 1 khách hàng mới sinh nhật vào tháng 7. Tự động thêm khách hàng này vào nhóm (có thể viết cron job)

50. Giao diện 'Tạo Voucher Khuyến Mại':
    Cần giao diện danh sách Voucher. Có nút "Thêm Voucher" tức là tạo lẻ 1 voucher theo điều kiện. Nút 'Tạo Hàng Loạt' là để tạo ra hàng loạt các voucher với các điều kiện áp dụng khác nhau, tùy theo báo cáo bán hàng + Báo cáo hoạt động của khách hàng

Giao diện Chi tiết 1 Voucher: đây là giao diện tạo Voucher Áp dụng cho từng khách hàng, hoặc áp dụng cho từng nhóm Khách hàng, áp dụng cho từng sản phẩm hoặc từng nhóm sản phẩm. Với nhiều kịch bản khác nhau. ví dụ: - Kịch bản 1: Tặng chính nó. Ví dụ: mua 10 chai dầu gội đầu Clear tặng 1 chai. - Kịch bản 2: Tặng sản phẩm khác. Ví dụ: Mua 5 hộp sữa bột tặng 1 bình pha sữa (Hàng Gift) - Kịch bản 3: Mua X tặng Y bậc thang. Ví dụ: Mua 10 tặng 1, mua 20 tặng 3, mua 50 tặng 10. (Chỉ cần lưu JSON dạng mảng array các mốc). - Kịch bản 4 - Mua chéo nhóm hàng: Mua 2 sản phẩm bất kỳ thuộc nhóm "Thực phẩm chức năng", tặng 1 hộp Vitamin C. - CHIẾT KHẤU THEO % (Percentage Discount) - Kịch bản 2.1 - Giảm trên tổng bill. Hóa đơn trên 1.000.000đ giảm 5% tổng giá trị thanh toán.
JSON Conditions: {"min_cart_value": 1000000}
JSON Rewards: {"reward_type": "discount_percent", "value": 5} - Kịch bản 2.2 - Giảm trên sản phẩm chỉ định: Mua 3 tuýp kem bôi da, được giảm 10% chỉ cho 3 tuýp đó (Các hàng khác trong bill giữ nguyên giá) - Kịch bản 2.3 - Happy Hour / Ngày vàng: Giảm 15% cho mọi hóa đơn xuất vào thứ 6 hàng tuần (Sử dụng kết hợp cột start_date và end_date).

    -Nhóm kịch bản: GIẢM TRỪ TIỀN MẶT (Fixed Amount)

        - Kịch bản 3.1 - Voucher giảm thẳng: Mua tổng toa thuốc trị giá 500.000đ, giảm ngay 50.000đ trừ thẳng vào hóa đơn.
            - JSON Conditions: {"min_cart_value": 500000}
            - JSON Rewards: {"reward_type": "fixed_discount", "value": 50000}

    - Kịch bản 3.2 - Trợ giá sản phẩm mới: Sản phẩm A giá gốc 120k, đợt này trợ giá bán đồng giá 99k/hộp (Giới hạn tối đa 3 hộp/khách).

    - Nhóm kịch bản: BÁN THEO COMBO (Bundle)
      - Kịch bản 4.1 - Combo cố định: Mua 1 Máy đo huyết áp + 1 Hộp que thử đường huyết. Tổng giá mua lẻ là 1.200.000đ, nhưng mua chung Combo chỉ còn 1.000.000đ.
        Cấu trúc dữ liệu:
        JSON Conditions: {"require_all": [{"product_id": 101, "qty": 1}, {"product_id": 102, "qty": 1}]}
        JSON Rewards: {"reward_type": "bundle_price", "value": 1000000}

    - Nhóm kịch bản: CÁ NHÂN HÓA KHÁCH HÀNG (Loyalty & Phân tệp) => Nếu thiết kế cột điều kiện khuyến mại là JSON (hoặc JSONB). có thể nhúng thêm điều kiện về customer một cách cực kỳ linh hoạt mà không cần sửa bảng Database:
      ví dụ:
      - Kịch bản 5.1 - Khách B2B vs B2C: Chương trình Mua 10 tặng 2 chỉ áp dụng cho Khách sỉ (B2B). Khách lẻ (B2C) không được áp dụng. Dữ liệu 'JSON Conditions: {"product_id": 4595, "min_qty": 10, "customer_type": "B2B"}'
      - Kịch bản 5.2 - Mừng sinh nhật: Khách hàng có ngày sinh nhật trong tháng hiện tại được tự động giảm 10% tổng bill.

51. Giao diện "Quản lý Chiến Dịch Marketing"
    Phần này cũng gồm 2 giao diện như các phần khác: Giao diện List danh sách các chiến dịch Marketing + Giao diện Chi tiết để CRUD chiến dịch MKT. 1. Giao diện List khá đơn giản vì nó là 1 giao diện 'Chỉ đọc' dữ liệu và hiển thị cho user. 2. Giao diện chi tiết 1 chiến dịch MKT: Đây là giao diện Phức tạp, để Tạo và Quản lý các chiến dịch Marketing. Có thể là 1 chiến dịch ngắn ngày (kéo dài trong 1 ngày) hoặc là một chiến dịch Dài hơi (Kéo dài 3 tháng) để nhằm đạt được 1 mục đích gì đó (ví dụ: nhắm mục tiêu tăng số lượng khách hàng nữ, độ tuổi từ 20 đến 50, số lượng từ 50 lên 150 người trong vòng 1 tháng )

        PHẦN 1: PHẦN MỀM MARKETING SẼ LÀM NHIỆM VỤ GÌ?
        - Tại các MNCs, phần mềm không chỉ là nơi lưu trữ, nó là một Người điều phối tự động (Campaign Orchestrator). Nó làm 3 nhiệm vụ chính:
          - Gắn kết Nguyên liệu: Ráp "Tệp khách hàng A" với "Voucher B" và "Thông điệp C".
          - Tự động hóa luồng phân phối (Workflow): Lên lịch gửi tin nhắn qua Zalo/App/SMS theo kịch bản có sẵn. Nó không gửi một cục, mà gửi theo luồng (Ví dụ: Gửi Zalo -> 2 ngày sau không đọc -> Gửi SMS -> Đọc rồi nhưng không mua -> Gửi App Push nhắc nhở).
          - A/B Testing: Tự động gửi 2 mẫu tin nhắn khác nhau cho 10% khách hàng đầu tiên, mẫu nào được click nhiều hơn sẽ tự động gửi cho 90% khách hàng còn lại.

        PHẦN 2: ĐIỀN DỮ LIỆU GÌ BAN ĐẦU KHI TẠO CHIẾN DỊCH?
        - Khi một Giám đốc Marketing (CMO) hoặc Quản lý tạo mới một chiến dịch trên phần mềm (Ví dụ: "Chiến dịch Xả hàng TPCN Tháng 6"), họ sẽ điền các Form dữ liệu đầu vào sau:
          - Thông tin cơ bản:
            - Tên chiến dịch (Campaign Name).
            - Mục tiêu chiến dịch (Objective): Tăng doanh thu, Xả hàng tồn, hay Tri ân khách VIP?
            - Thời gian chạy: Start Date & End Date.

          - Ngân sách (Budgeting):
            - Tổng ngân sách cho phép (Ví dụ: 50.000.000 VNĐ). Hệ thống sẽ tự động dừng chiến dịch khi tổng số tiền Voucher được sử dụng + Phí gửi Zalo/SMS chạm mức này.

          - Nhắm mục tiêu (Targeting):
            - Chọn "Nhóm khách hàng" (Tích chọn từ Module 49 ở trên).

          - Ưu đãi (Offering):
            - Chọn "Voucher" (Tích chọn từ Module 50).

          - Kênh & Thông điệp (Channels & Content):
            - Đây là phần thiết lập như 1 danh sách (table) 'Điểm chạm (Touchpoints)' vào khách hàng, được gắn thời gian. Khu vực này có nút bấm 'Thêm Điểm Chạm' (thực tế là thêm 1 hành động phải làm và gửi đến khách hàng). Ví dụ khi bấm 'Thêm điểm Chạm', user phải chọn: -
              - Chọn thời gian:
                - Thời gian Tuyệt đối (Absolute Scheduled Time):
                  - Ứng dụng: Dành cho các chiến dịch có ngày giờ cụ thể, ví dụ như Flash Sale, Lễ tết, hoặc bài đăng Facebook.
                    Thiết lập: Chọn chính xác Ngày - Giờ - Phút để hệ thống tự nhả đạn.
                    Ví dụ trong hệ thống:
                    Hành động 1: 08:00 ngày 01/06/2026 -> Đăng bài lên Fanpage khởi động chiến dịch.
                    Hành động 2: 09:00 ngày 15/06/2026 -> Bắn Zalo nhắc nhở Sale giữa tháng.

                - Thời gian Tương đối / Độ trễ (Relative Delay Time)
                  - Ứng dụng: Dành cho các chuỗi tin nhắn chăm sóc (Drip Campaign) nối tiếp nhau, tính toán dựa trên ngày khách hàng bắt đầu lọt vào chiến dịch.
                  - Thiết lập: Gửi sau X ngày / Y giờ kể từ sự kiện trước đó.
                  - Ví dụ trong hệ thống:
                    Ngày 0: Khách hàng (Quầy thuốc) bấm đăng ký tài khoản B2B thành công -> Hệ thống tự động gửi Zalo ZNS chào mừng.
                    Ngày +3: Tự động gửi Zalo tặng Voucher giảm giá đơn đầu tiên.

                - Thời gian theo Hành vi (Behavioral Trigger Time)
                  - Ứng dụng: Kích hoạt ngay khi khách hàng có một "động thái" cụ thể, đánh trúng tâm lý lúc họ đang có nhu cầu cao nhất.
                  - Thiết lập: Nếu Khách hàng thực hiện Hành vi A, chờ X giờ, thực hiện Hành động B.
                  - Ví dụ trong hệ thống:
                    Khách hàng cho mặt hàng "Sữa tăng cân" vào giỏ nhưng thoát App không thanh toán.
                    Hệ thống đếm ngược đúng 24 giờ sau, tự động bắn thông báo Push App: "Anh chị bỏ quên giỏ hàng kìa, Nam Việt tặng thêm anh chị mã Freeship nhé!"
                    Ngày +7 (Kể từ ngày gửi Voucher): Nếu khách chưa mua, tự động gửi SMS nhắc hạn sử dụng Voucher sắp hết. - Chọn kênh gửi: Facebook, Zalo OA SMS, Web Pop-up, hay Mobile App Push. .v.v.v - Soạn nội dung tin nhắn / Bài viết (Có thể chèn biến cá nhân hóa: "Chào [Tên Khách Hàng], Nam Việt tặng anh/chị...").

          - Lựa chọn Kích hoạt:
            - Chiến dịch Chạy 1 Lần (Scheduled/Batch) : cứ so sánh thời điểm hiện tại với ngày bắt đầu và ngày kết thúc để tự On / Off các hành động.
            - Chiến dịch Tự Động Hóa (Trigger-based/Always-On) : đây là chuỗi các hành động tự động cập nhật và chạy. Ví dụ: Chiến dịch "Chúc mừng sinh nhật". User thiết lập chuỗi: Gửi tin nhắn chúc mừng (Ngày 0) -> Nhắc lại mã giảm giá (Ngày +2). Chiến dịch này sẽ được bật ở trạng thái "Đang chạy", Cứ hôm nay hệ thống quét thấy ai đến sinh nhật -> Bốc người đó ném vào Ngày 0 của chuỗi. Ngày mai có người khác sinh nhật -> Lại bốc người đó ném vào Ngày 0.

    - PHẦN 3: THEO DÕI CÁC DỮ LIỆU GÌ KHI CHIẾN DỊCH ĐANG CHẠY?
      Đây là phần "ăn tiền" nhất của các MNCs. Họ không đo lường bằng cảm giác, họ nhìn vào Phễu chuyển đổi chiến dịch (Campaign Funnel). Phần mềm cần có một Dashboard báo cáo Real-time (Thời gian thực) cho từng chiến dịch với các chỉ số sau:
      1. Nhóm chỉ số Phân phối & Tương tác (Engagement Metrics):
      - Sent / Delivered: Đã gửi thành công bao nhiêu tin nhắn Zalo/App.
      - Open Rate (Tỷ lệ mở): Bao nhiêu % khách hàng đã mở tin nhắn đọc.
      - Click-Through Rate (CTR): Bao nhiêu % khách hàng bấm vào nút "Xem ngay" hoặc "Nhận Voucher" trong tin nhắn.
      2. Nhóm chỉ số Khuyến mãi (Voucher Metrics):
      - Claimed (Đã lưu): Bao nhiêu khách đã lưu Voucher vào ví của họ (Nhưng chưa mua).
      - Redeemed (Đã sử dụng): Bao nhiêu khách thực sự đã áp mã Voucher đó để thanh toán thành công. (Đây là chỉ số quan trọng nhất để đánh giá Voucher đó có đủ hấp dẫn hay không).
      3. Nhóm chỉ số Hiệu quả Kinh doanh (Business Metrics):
      - Total Campaign Revenue (Doanh thu chiến dịch): Tổng giá trị các đơn hàng có sử dụng Voucher của chiến dịch này.
      - Campaign ROI (Tỷ suất hoàn vốn): Doanh thu thu về trừ đi (Tiền giảm giá Voucher + Tiền phí gửi tin nhắn Zalo/SMS).
      - CPA (Cost Per Acquisition): Trung bình công ty tốn bao nhiêu tiền (chi phí gửi tin + tiền giảm giá) để có được 1 đơn hàng thành công trong chiến dịch này.

52. Giao diện 'Khảo Sát Khách Hàng'

- Cũng gồm 2 giao diện: List và Tạo khảo sát chi tiết.
- Dùng lại các phần Khởi tạo của bên marketing, nhưng sẽ khác phần form 'PHẦN 2: ĐIỀN DỮ LIỆU GÌ BAN ĐẦU KHI TẠO CHIẾN DỊCH?' thay vì là tin nhắn thì đây là 1 form để Khách hàng trả lời câu hỏi. Và các câu trả lời này được lưu JSON vào Database của khách hàng.

53. Giao diện "Quản lý Bài Viết và Nội dung Số":

- Giao diện tạo Bài viết và Nội dung khác nhau trên các kênh khác nhau, ví dụ: Facebook; videos ...
- Vẫn gồm 2 giao diện List và Tạo Bài Viết Chi tiết. Ở Mục Quản lý Chiến Dịch Marketing cần liên kết được với 1 hoặc nhiều bài viết này.
- Cần lấy được dữ liệu báo cáo từ các bài viết hoặc nội dung này. Ví dụ: lượt xem; lượt tương tác; ...
- Quản lý các thông tin:
  - Đăng kênh nào? (web sỉ/web lẻ; mobile app;... )
  - Bao giờ Đăng? (Quản lý lịch đăng)
  - Giới hạn nhóm khách nào xem không?

## NHÓM CÁC GIAO DIỆN VẬN HÀNH:

54. Mua Hàng & Nhập Kho: quản lý Vòng đời của 1 đơn mua hàng và nhập vào kho + Thanh toán: Tạo đơn Nháp => Đã đặt & Chờ giao => Nhận hàng & Nhập vào Kho => Thanh toán NCC => Quản lý Hóa Đơn VAT.
    Yêu cầu các tính năng của giao diện này:

- Tại giao diện List: có tính năng "Tạo đơn Mua Hàng" hàng loạt dựa vào logic: Tồn kho Thực tế của sản phẩm tại kho được chỉ định nhỏ hơn Tồn Min Của Sản Phẩm (được cài đặt sẵn) sẽ tạo đơn mua hàng với số lượng Tồn Max của sản phẩm tại kho trừ đi Tồn Thực Tế. Phiếu Mua hàng được gom theo Nhà Cung Cấp.
- Tại giao diện Chi tiết 1 đơn Mua Hàng cần có các tính năng
  - Upload phiếu xuất kho / Hóa đơn VAT của NCC => tự động điền Lô Hạn sử dụng vào phiếu mua hàng
  - Sau khi Upload Phiếu Xuất Kho / Hóa đơn VAT của NCC => Phân biệt được đâu là hàng tặng, đâu là hàng mua => có nút "Tạo Chương trình Bán Hàng" giống như của Nhà Cung Cấp => Ghi vào Database chương trình bán hàng / Khuyến mại sau khi user bấm "Đồng Ý" (tức là xác nhận tạo chương trình bán hàng)
  - Tự động Lựa chọn các chương trình của Nhà Cung Cấp, Dựa vào các chương trình + Hợp Đồng Của Nhà Cung Cấp => tính ra giá vốn Nhập Vào của sản phẩm tại Lô được nhập lưu vào Database để tính ra Lãi Lỗ trong báo cáo tài chính.
  - Có nút "Gửi Thông Báo Nội Bộ" tức là thông báo cho nội bộ (Ví dụ: dược sĩ / bác sĩ ...) biết rằng có sản phẩm mới về.
  - Gắn đơn mua hàng này, với 1 hóa đơn VAT cụ thể, nếu đơn mua hàng này đã được thanh toán thì cần cập nhật số tiền đã thanh toán cho hóa đơn VAT đó.

55. Xuất Kho & Giao Hàng: Quản lý vòng đời 1 đơn Bán Hàng (hoặc chuyển hàng) + Giao Hàng + Thanh Toán.
    bao gồm: Tạo Đơn Hàng/Phiếu Xuất kho => Kho Đóng Gói + Xuất Hóa Đơn VAT => Gửi hàng + Bộ giấy tờ cho Vận Chuyển => Khách hàng xác nhận đã nhận hàng => Thanh toán và cập nhật trạng thái đơn hàng
    - Cần có thêm tính năng cho trường hợp Nhân viên Giao Hàng Nội Bộ, tại thời điểm giao hàng cho khách hàng, nếu khách hàng 'Thanh toán tiền mặt' => NV Giao Hàng bấm nút 'Đã Thanh Toán Tiền Mặt' trên đơn hàng => Đơn hàng được cập nhật 'Đã Thanh Toán' => nợ của Khách hàng giảm => Nhưng số tiền này chưa được thu vào Quỹ, và Nhân viên Giao Hàng cần nộp lại tiền vào Quỹ thông qua việc thủ Quỹ bấm nút 'Xác Nhận đã Thu'.
    - Khi khách hàng thanh toán 1 cục (Ví dụ: 10 triệu, bằng bất cứ hình thức thu tiền nào), thì bảng finance_transactions vẫn hiển thị số tiền thu của 1 cục đó. Nhưng sẽ có thêm icon expan để mở rộng ra và hiển thị rằng 1 cục tiền đó phân bổ chi tiết số tiền cho từng đơn hàng nào (ví dụ: phân bổ đơn hàng A1 2 triệu; phân bổ đơn hàng A2 3 triệu; phân bổ đơn hàng A3 5 triệu) Yêu cầu, Đơn hàng xa nhất sẽ được cập nhật thanh toán trước (hoặc do cài đặt của Admin trong phần cấu hình hệ thống) [ ví dụ: có thể tạo bảng finance_transaction_allocations để gắn 1 cục tiền với nhiều đơn hàng ]

56. Kiểm Kê: tạo phiếu kiểm tại mỗi Kho.
    - Tạo kiểm kê Hàng Hóa (hàng bán)
    - Tạo kiểm kê Quà Tặng
    - Tạo kiểm kê Tài Sản

## NHÓM CÁC GIAO DIỆN KINH DOANH:

57. Giao diện "Tạo Đơn Khách Sỉ"

- ngoài các tính năng cốt lõi, cần thêm các tính năng:
- Tìm kiếm khách hàng + Cập nhật (sửa) thông tin khách hàng thông qua popup hoặc drawer ngay tại giao diện tạo đơn hàng.
  Khi Giỏ hàng thay đổi, Backend không chỉ tính giá, mà còn trả về Array upsell_suggestions (VD: "Mua thêm 200k nữa để được giảm 10%").
- Có nút "Gửi Qua hãng Vận Chuyển" => bấm nút này, 1 đơn giao vận sẽ được tạo thông qua các hãng vận chuyển thứ 3 (Ví dụ: Shopee; Giao hàng Nhanh; Viettel Post; ...)

58. Giao diện "Danh sách đơn hàng Sỉ"
59. Giao diện "Chăm Sóc Khách Hàng Sỉ"

60. Giao diện "Tạo đơn hàng Khách Lẻ"
    Ngoài các chức năng cốt lõi, cần thêm các tính năng:

- Cần có chức năng: Chọn bệnh => Lấy tuổi của Khách hàng hiện tại (Để phân biệt người lớn / trẻ nhỏ ...) => Lấy đơn thuốc mẫu, và tự động điền vào phần giỏ hàng trên giao diện, tự điền thông tin Hướng Dẫn Sử Dụng từ đơn thuốc mẫu.

61. Giao diện "Danh sách Đơn Khách Lẻ"
62. Giao diện "Chăm sóc Khách Lẻ"

- Tự động sinh ra từ các đơn mua hàng / các dịch vụ mà khách lẻ đã sử dụng => có thông tin Bệnh chi tiết để Dược Sĩ / Lễ Tân gọi điện hoặc liên hệ lại khách hàng.

63. Giao diện "Tạo lịch Hẹn" => Giao diện này dành cho Điều Dưỡng/Lễ Tân sử dụng. Đây cũng là Giao diện để bán các dịch vụ luôn. Ví dụ: Bán Mũi tiêm Vắc xin; Bán gói tiêm Vắc Xin; Bán gói Siêu Âm; Bán dịch vụ 'Vệ Sinh và Khâu Vết thương'; ...
64. Giao diện "Khám Bệnh": Dành cho Bác Sĩ. Chỉ khi khách có được đặt lịch (Đặt trước hoặc ngay lập tức), bác sĩ sẽ nhìn thấy khách ở hàng đợi và thực hiện 'Khám Bệnh'. Sau đó chỉ định Bệnh nhân thực hiện các Thủ Thuật hoặc Tiêm.
65. Giao diện "Tiêm & Thủ Thuật": dành cho Điều Dưỡng thực hiện Tiêm Vắc Xin hoặc Thực hiện các thủ thuật dịch vụ (Ví dụ: Vệ sinh và Khâu Vết thương)
66. Giao diện "Cận Lâm Sàng": Giao diện dành cho các kỹ thuật viên, thực hiện các hành động. Ví dụ: Làm xét nghiệm; Thực hiện Siêu Âm; Thực hiện Chụp X Quang; ...

## NHÓM CÁC GIAO DIỆN TÀI CHÍNH KẾ TOÁN

67. Quản lý Thu / Chi (Finance Transaction)
68. Giao diện "Đối Soát Giao Dịch"
69. Giao diện 'Quản lý Hóa Đơn VAT Mua Vào'
70. Giao diện 'Quản lý Hóa đơn VAT Bán Ra'
71. Giao diện 'Sổ Cái Kế Toán': là giao diện hiển thị các Bút toán đã được Hạch Toán.

## NHÓM CÁC GIAO DIỆN BÁO CÁO:

81. Báo cáo Khách Hàng và Chăm Sóc Khách Hàng: thể hiện chi tiết các chỉ số của khách hàng, ví dụ như: Tổng bao nhiêu khách hàng; Bao nhiêu Khách quay lại; Lần mua cuối; Bao nhiêu khách hàng được Gọi điện Chăm sóc hỏi thăm; Nhu cầu của khách hàng (dựa vào Khảo Sát); v.v.
82. Báo Cáo Nhân Sự: báo cáo về chấm công; hiệu suất của Nhân Sự (thông qua hoàn thành KPIs, Hoàn thành các công việc được giao, ...)
83. Báo Cáo Kinh Doanh và Lợi Nhuận: hiệu quả kinh doanh và Lợi Nhuận
84. Báo cáo Vận Hành: Báo cáo Kho và Vận hành, giao vận; các sản phẩm cận date; vòng đời sản phẩm tồn kho; lợi nhuận dựa trên mỗi sản phẩm; ...
85. Báo cáo Marketing và Truyền thông: Tỉ suất sinh lợi nhuận từ các chiến dịch; Các chỉ số cần theo dõi tại phần mục Marketing; ...

## NHÓM CÁC GIAO DIỆN CẤU HÌNH AI và các QUYỀN mà AI được làm
