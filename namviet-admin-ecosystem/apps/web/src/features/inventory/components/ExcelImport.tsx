import React, { useState } from 'react';
import { Upload, Button, Table, Alert, Card, Select, Form, Input } from 'antd';
import { UploadOutlined, FileExcelOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import * as XLSX from 'xlsx';
import { useProducts } from '@/features/products/hooks/useProducts';
import { useBatches } from '@/features/batches/hooks';
import { useCreateTransaction } from '../hooks';
import { useWarehouses } from '@/features/warehouses/hooks';
import { toast } from 'sonner';

interface ExcelRow {
  sku: string;
  batch_code: string;
  quantity: number;
  price: number;
  mfg_date?: string;
  exp_date?: string;
  [key: string]: any;
}

interface ParsedItem {
  key: string;
  sku: string;
  product_id?: number;
  product_name?: string;
  batch_code: string;
  batch_id?: number;
  quantity: number;
  price: number;
  unit_id?: number;
  isValid: boolean;
  errorMsg?: string;
}

export const ExcelImport = () => {
  const [fileList, setFileList] = useState<any[]>([]);
  const [parsedData, setParsedData] = useState<ParsedItem[]>([]);
  const [warehouseId, setWarehouseId] = useState<number>();
  const [notes, setNotes] = useState('');
  
  const { data: products = [] } = useProducts();
  const { data: batches = [] } = useBatches();
  const { data: warehouses = [] } = useWarehouses();
  const createMutation = useCreateTransaction();

  const handleFileUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const data = e.target?.result;
        const workbook = XLSX.read(data, { type: 'binary' });
        const firstSheetName = workbook.SheetNames[0];
        const worksheet = workbook.Sheets[firstSheetName];
        
        // Cấu trúc cột kỳ vọng: SKU, Lô, Số lượng, Đơn giá
        const rawData = XLSX.utils.sheet_to_json<any>(worksheet);
        
        const validatedData: ParsedItem[] = rawData.map((row, index) => {
          const sku = String(row['SKU'] || '').trim();
          const batch_code = String(row['Mã Lô'] || row['Lô'] || '').trim();
          const quantity = Number(row['Số lượng'] || row['SL']);
          const price = Number(row['Đơn giá'] || row['Giá']);

          let isValid = true;
          const errorMsgs = [];

          if (!sku) {
            isValid = false;
            errorMsgs.push('Thiếu SKU');
          }
          if (!batch_code) {
            isValid = false;
            errorMsgs.push('Thiếu Mã Lô');
          }
          if (isNaN(quantity) || quantity <= 0) {
            isValid = false;
            errorMsgs.push('Số lượng không hợp lệ');
          }

          // Kiểm tra xem SKU có tồn tại trong hệ thống không
          const product = products.find(p => p.sku === sku);
          if (sku && !product) {
            isValid = false;
            errorMsgs.push('SKU không tồn tại');
          }

          // Tìm ID của Lô (Batch), nếu không có thì hệ thống backend phải hỗ trợ tự tạo,
          // nhưng ở đây giả sử Batch đã được tạo trước, hoặc backend sẽ xử lý. 
          // Tạm thời ta chỉ pass batch_id nếu tìm thấy.
          const batch = batches.find(b => b.batch_code === batch_code && b.product_id === product?.id);

          return {
            key: `row-${index}`,
            sku,
            product_id: product?.id,
            product_name: product?.name,
            batch_code,
            batch_id: batch?.id,
            quantity,
            price: isNaN(price) ? 0 : price,
            isValid,
            errorMsg: errorMsgs.join(', '),
          };
        });

        setParsedData(validatedData);
        toast.success(`Đã đọc ${validatedData.length} dòng dữ liệu`);
      } catch (error) {
        toast.error('Lỗi đọc file Excel. Vui lòng kiểm tra định dạng.');
        console.error(error);
      }
    };
    reader.readAsBinaryString(file);
    return false; // Prevent automatic upload
  };

  const handleDownloadTemplate = () => {
    const ws = XLSX.utils.json_to_sheet([
      { 'SKU': 'SP001', 'Mã Lô': 'L001', 'Số lượng': 100, 'Đơn giá': 50000 },
      { 'SKU': 'SP002', 'Mã Lô': 'L002', 'Số lượng': 50, 'Đơn giá': 120000 }
    ]);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, "Nhập Kho");
    XLSX.writeFile(wb, "Template_NhapKho.xlsx");
  };

  const handleSubmit = () => {
    if (!warehouseId) {
      toast.error('Vui lòng chọn Chi nhánh nhập kho');
      return;
    }

    const validItems = parsedData.filter(item => item.isValid);
    if (validItems.length === 0) {
      toast.error('Không có dữ liệu hợp lệ để nhập kho');
      return;
    }

    const payload = {
      warehouse_id: warehouseId,
      transaction_type: 'IMPORT',
      reference_type: 'EXCEL',
      notes: notes || 'Nhập kho bằng file Excel',
      items: validItems.map(item => ({
        product_id: item.product_id!,
        batch_id: item.batch_id || 0, // Backend cần hỗ trợ tạo Lô nếu truyền 0 + batch_code
        quantity: item.quantity,
        price: item.price,
      }))
    };

    createMutation.mutate(payload as any, {
      onSuccess: () => {
        toast.success(`Nhập kho thành công ${validItems.length} sản phẩm`);
        setParsedData([]);
        setFileList([]);
        setNotes('');
      },
      onError: (err) => toast.error(`Lỗi: ${err.message}`)
    });
  };

  const columns = [
    {
      title: 'Trạng thái',
      key: 'status',
      width: 100,
      render: (_: any, record: ParsedItem) => (
        record.isValid ? 
          <div className="flex items-center text-green-600 gap-1"><CheckCircleOutlined className="w-4 h-4"/> Hợp lệ</div> : 
          <div className="flex items-center text-red-500 gap-1"><CloseCircleOutlined className="w-4 h-4"/> Lỗi</div>
      )
    },
    { title: 'SKU', dataIndex: 'sku', key: 'sku', width: 120 },
    { 
      title: 'Sản phẩm', 
      key: 'product',
      render: (_: any, record: ParsedItem) => record.product_name || <span className="text-red-400 italic">Không tìm thấy</span>
    },
    { title: 'Mã Lô', dataIndex: 'batch_code', key: 'batch', width: 120 },
    { title: 'Số lượng', dataIndex: 'quantity', key: 'quantity', width: 100, align: 'right' as const },
    { 
      title: 'Đơn giá', 
      dataIndex: 'price', 
      key: 'price', 
      width: 120, 
      align: 'right' as const,
      render: (val: number) => val ? val.toLocaleString('vi-VN') : '0'
    },
    { 
      title: 'Chi tiết lỗi', 
      dataIndex: 'errorMsg', 
      key: 'errorMsg',
      render: (text: string) => <span className="text-red-500 text-xs">{text}</span>
    },
  ];

  const totalValid = parsedData.filter(p => p.isValid).length;
  const totalInvalid = parsedData.length - totalValid;

  return (
    <div className="space-y-6">
      <Card title={<span className="flex items-center gap-2"><FileExcelOutlined /> Import Dữ liệu Nhập Kho</span>}>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div>
            <h3 className="font-semibold mb-2">1. Tải File Mẫu</h3>
            <p className="text-sm text-gray-500 mb-4">Sử dụng file Excel mẫu để đảm bảo cấu trúc dữ liệu chính xác (SKU, Mã Lô, Số lượng, Đơn giá).</p>
            <Button onClick={handleDownloadTemplate}>Tải Template_NhapKho.xlsx</Button>
          </div>

          <div>
            <h3 className="font-semibold mb-2">2. Cấu hình & Upload</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Chi nhánh nhập hàng <span className="text-red-500">*</span></label>
                <Select 
                  className="w-full"
                  placeholder="Chọn chi nhánh"
                  value={warehouseId}
                  onChange={setWarehouseId}
                  options={warehouses.map(w => ({ value: w.id, label: w.name }))}
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-1">Ghi chú</label>
                <Input.TextArea 
                  rows={2} 
                  placeholder="VD: Nhập hàng đợt 1 tháng 10..."
                  value={notes}
                  onChange={e => setNotes(e.target.value)}
                />
              </div>

              <div>
                <Upload
                  accept=".xlsx, .xls"
                  beforeUpload={handleFileUpload}
                  fileList={fileList}
                  onChange={({ fileList }) => setFileList(fileList)}
                  maxCount={1}
                >
                  <Button icon={<UploadOutlined />} type="primary" ghost>Chọn File Excel</Button>
                </Upload>
              </div>
            </div>
          </div>
        </div>
      </Card>

      {parsedData.length > 0 && (
        <Card title="Dữ liệu xem trước">
          <div className="flex gap-4 mb-4">
            <Alert type="info" message={`Tổng cộng: ${parsedData.length} dòng`} showIcon />
            <Alert type="success" message={`Hợp lệ: ${totalValid} dòng`} showIcon />
            {totalInvalid > 0 && <Alert type="error" message={`Lỗi: ${totalInvalid} dòng`} showIcon />}
          </div>

          <Table 
            columns={columns} 
            dataSource={parsedData} 
            size="small" 
            pagination={{ pageSize: 50 }}
            scroll={{ y: 400 }}
            rowClassName={(record) => !record.isValid ? 'bg-red-50' : ''}
          />

          <div className="flex justify-end mt-4">
            <Button 
              type="primary" 
              size="large" 
              className="bg-blue-600"
              disabled={totalValid === 0 || !warehouseId || createMutation.isPending}
              loading={createMutation.isPending}
              onClick={handleSubmit}
            >
              Xác nhận Nhập kho ({totalValid} hợp lệ)
            </Button>
          </div>
        </Card>
      )}
    </div>
  );
};
