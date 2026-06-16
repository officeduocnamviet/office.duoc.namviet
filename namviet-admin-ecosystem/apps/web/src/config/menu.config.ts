import { 
  Store, Activity, PackageSearch, Boxes, Receipt, Headphones,
  Users, Megaphone, BarChart3, Settings, MessageSquare, BrainCircuit
} from 'lucide-react';

export const MENU_DATA = [
  {
    groupTitle: 'TRẠM LÀM VIỆC',
    items: [
      {
        id: 'retail', // Dùng làm key cho accordion nếu có children
        name: 'Cửa hàng (Bán lẻ)',
        icon: Store,
        children: [
          { href: '/retail/create', name: 'Tạo đơn mới' },
          { href: '/retail/list', name: 'Danh sách Đơn hàng' },
          { href: '/retail/cskh', name: 'Chăm sóc khách hàng' },
        ]
      },
      {
        id: 'customers',
        name: 'Khách hàng & Đối tác',
        icon: Users,
        children: [
          { href: '/customers', name: 'Quản lý Khách hàng' },
        ]
      },
      {
        id: 'catalog',
        name: 'Sản phẩm & Danh mục',
        icon: Boxes,
        children: [
          { href: '/products', name: 'Quản lý Sản phẩm' },
          { href: '/categories', name: 'Nhóm Sản phẩm' },
          { href: '/manufacturers', name: 'Nhà sản xuất' },
        ]
      },
      {
        id: 'clinical',
        name: 'Y Tế & Lâm sàng',
        icon: Activity,
        children: [
          { href: '/clinical/appointments', name: 'Lịch hẹn (Appointments)' },
          { href: '/clinical/queues', name: 'Hàng đợi (Queues)' },
          { href: '/clinical/visits', name: 'Hồ sơ khám (Visits)' },
        ]
      },
      {
        id: 'b2b',
        name: 'Bán Sỉ (B2B)',
        icon: PackageSearch,
        children: [
          { href: '/b2b/create', name: 'Tạo đơn mới' },
          { href: '/b2b/list', name: 'Danh sách Đơn hàng' },
          { href: '/b2b/cskh', name: 'Chăm sóc khách hàng (B2B)' },
        ]
      },
      {
        id: 'inventory',
        name: 'Vận hành Kho',
        icon: Boxes,
        children: [
          { href: '/inventory/import', name: 'Nhập hàng' },
          { href: '/inventory/export', name: 'Xuất hàng' },
          { href: '/inventory/transfer', name: 'Chuyển hàng (Kho)' },
          { href: '/inventory/check', name: 'Kiểm kê' },
          { href: '/inventory/batches', name: 'Quản lý Lô Hàng' },
          { href: '/inventory/gifts', name: 'Quản lý Quà tặng' },
          { href: '/inventory/assets', name: 'Quản lý Vật tư (Tài sản)' },
          { href: '/inventory/shipping', name: 'Giao vận (Theo dõi)' },
        ]
      },
      {
        id: 'finance',
        name: 'Tài chính & Kế toán',
        icon: Receipt,
        children: [
          { href: '/finance/funds', name: 'Quỹ & Tài khoản' },
          { href: '/finance/transactions', name: 'Giao dịch Tài chính' },
          { href: '/finance/coa', name: 'Hệ thống Tài khoản' },
          { href: '/finance/journals', name: 'Nhật ký chung' },
        ]
      },
      {
        id: 'support',
        name: 'Hỗ trợ & Tương tác',
        icon: Headphones,
        children: [
          { href: '/support/bot', name: 'Tiếp nhận Chatbot' },
          { href: '/support/fb', name: 'Kênh Facebook' },
          { href: '/support/zalo', name: 'Kênh Zalo' },
          { href: '/support/tiktok', name: 'Kênh TMĐT & TikTok' },
        ]
      }
    ]
  },
  {
    groupTitle: 'PHÒNG ĐIỀU HÀNH',
    items: [
      {
        id: 'hr',
        name: 'Nhân sự & Tính lương',
        icon: Users,
        children: [
          { href: '/hr/employees', name: 'Hồ sơ nhân viên' },
          { href: '/hr/contracts', name: 'Hợp đồng lao động' },
          { href: '/hr/attendance', name: 'Chấm công' },
          { href: '/hr/shifts', name: 'Lịch phân ca' },
          { href: '/hr/handovers', name: 'Bàn giao ca' },
          { href: '/hr/payrolls', name: 'Bảng lương' },
          { href: '/hr/training', name: 'Khóa đào tạo' },
        ]
      },
      {
        id: 'marketing',
        name: 'Marketing & PTKH',
        icon: Megaphone,
        children: [
          { href: '/marketing/campaigns', name: 'Chiến dịch Marketing' },
          { href: '/marketing/segment', name: 'Nhóm khách hàng' },
          { href: '/marketing/voucher', name: 'Voucher & CTKM' },
          { href: '/marketing/sms', name: 'Tin nhắn hàng loạt' },
        ]
      },
      {
        id: 'reports',
        name: 'Báo cáo tổng hợp',
        icon: BarChart3,
        children: [
          { href: '/reports/business', name: 'Báo cáo Kinh doanh' },
          { href: '/reports/finance', name: 'Báo cáo Tài chính' },
          { href: '/reports/inventory', name: 'Báo cáo Kho' },
          { href: '/reports/marketing', name: 'Báo cáo Marketing' },
        ]
      },
      {
        id: 'system',
        name: 'Hệ thống & Tích hợp',
        icon: Settings,
        children: [
          { href: '/system/configs', name: 'Cấu hình chung' },
          { href: '/system/approvals', name: 'Trung tâm Phê duyệt' },
          { href: '/system/integrations', name: 'Đối tác & Webhook' },
          { href: '/companies', name: 'Quản lý Công ty' },
          { href: '/warehouses', name: 'Quản lý Chi nhánh' },
          { href: '/roles', name: 'Quản lý Phân quyền' },
          { href: '/users', name: 'Quản lý Nhân sự' },
        ]
      },
      {
        id: 'communications',
        name: 'Giao tiếp Nội bộ',
        icon: MessageSquare,
        children: [
          { href: '/communications/chat', name: 'Kênh Chat Nội bộ' },
        ]
      },
      {
        id: 'ai_ecosystem',
        name: 'Hệ sinh thái AI',
        icon: BrainCircuit,
        children: [
          { href: '/ai/knowledge', name: 'Tri thức AI (Vectors)' },
          { href: '/ai/supervisor', name: 'Giám sát Chatbot' },
        ]
      }
    ]
  }
];

export const getPageTitleByHref = (currentHref: string) => {
  if (currentHref === '/' || currentHref === '/dashboard') return 'Dashboard';
  
  for (const group of MENU_DATA) {
    for (const item of group.items) {
      if ('href' in item && item.href === currentHref) return item.name;
      if (item.children) {
        const child = item.children.find(c => currentHref.startsWith(c.href));
        if (child) return child.name;
      }
    }
  }
  return 'Trang tính năng';
};
