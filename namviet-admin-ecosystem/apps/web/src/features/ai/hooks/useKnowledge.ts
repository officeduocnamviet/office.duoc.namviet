import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { knowledgeApi, CreateMedicalVectorReq, CreateProductVectorReq } from '../api/knowledgeApi';
import { toast } from 'sonner';

export const useMedicalVectors = () => {
  return useQuery({
    queryKey: ['medical_vectors'],
    queryFn: knowledgeApi.getMedicalVectors,
  });
};

export const useCreateMedicalVector = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateMedicalVectorReq) => knowledgeApi.createMedicalVector(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medical_vectors'] });
      toast.success('Thêm tri thức Y tế thành công');
    },
    onError: (err: any) => toast.error(err.response?.data?.error || 'Có lỗi xảy ra')
  });
};

export const useProductVectors = () => {
  return useQuery({
    queryKey: ['product_vectors'],
    queryFn: knowledgeApi.getProductVectors,
  });
};

export const useCreateProductVector = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateProductVectorReq) => knowledgeApi.createProductVector(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['product_vectors'] });
      toast.success('Thêm Vector Sản phẩm thành công');
    },
    onError: (err: any) => toast.error(err.response?.data?.error || 'Có lỗi xảy ra')
  });
};
