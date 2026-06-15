import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { productUnitApi } from '../api/unitApi';
import { 
  ProductUnitsCreateProductUnitRequest as CreateProductUnitRequest, 
  ProductUnitsUpdateProductUnitRequest as UpdateProductUnitRequest 
} from '@namviet/shared-types/src/backend.d';

export const unitKeys = {
  all: (productId: number) => ['productUnits', productId] as const,
};

export const useProductUnits = (productId: number) => {
  return useQuery({
    queryKey: unitKeys.all(productId),
    queryFn: () => productUnitApi.getByProductId(productId),
  });
};

export const useCreateProductUnit = (productId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateProductUnitRequest) => productUnitApi.create(productId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: unitKeys.all(productId) });
    },
  });
};

export const useUpdateProductUnit = (productId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ unitId, data }: { unitId: number; data: UpdateProductUnitRequest }) => 
      productUnitApi.update(productId, unitId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: unitKeys.all(productId) });
    },
  });
};

export const useDeleteProductUnit = (productId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (unitId: number) => productUnitApi.delete(productId, unitId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: unitKeys.all(productId) });
    },
  });
};
