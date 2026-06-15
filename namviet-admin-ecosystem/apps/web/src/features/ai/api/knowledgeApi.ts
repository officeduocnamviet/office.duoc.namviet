import { apiClient } from '@/lib/axios';

export type MedicalVector = import('@namviet/shared-types/src/backend.d').KnowledgeVectorsMedicalKnowledgeVector;
export type CreateMedicalVectorReq = import('@namviet/shared-types/src/backend.d').KnowledgeVectorsCreateMedicalKnowledgeVectorRequest;

export type ProductVector = import('@namviet/shared-types/src/backend.d').KnowledgeVectorsProductVector;
export type CreateProductVectorReq = import('@namviet/shared-types/src/backend.d').KnowledgeVectorsCreateProductVectorRequest;

export const knowledgeApi = {
  getMedicalVectors: async (): Promise<MedicalVector[]> => {
    const { data } = await apiClient.get<MedicalVector[]>('/api/medical-vectors');
    return data;
  },

  createMedicalVector: async (req: CreateMedicalVectorReq): Promise<MedicalVector> => {
    const { data } = await apiClient.post<MedicalVector>('/api/medical-vectors', req);
    return data;
  },

  getProductVectors: async (): Promise<ProductVector[]> => {
    const { data } = await apiClient.get<ProductVector[]>('/api/product-vectors');
    return data;
  },

  createProductVector: async (req: CreateProductVectorReq): Promise<ProductVector> => {
    const { data } = await apiClient.post<ProductVector>('/api/product-vectors', req);
    return data;
  }
};
