import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../api-client.ts';
import { toast } from 'react-toastify';
import { aliexpressQueryKeys } from './aliexpress-query-key.ts';
import { ProductTypeSource } from '../../types/product.ts';

const crawlAliProductsFn = async (aliexpressIds: string[]) => {
  const body = {
    productIds: aliexpressIds,
    source: ProductTypeSource.Aliexpress,
  };

  const response = await apiClient.post('product/crawl-aliexpress', body);
  return response.data;
};

export function useCrawlAliProducts() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: crawlAliProductsFn,
    onMutate: async () => {
      await queryClient.cancelQueries({
        queryKey: aliexpressQueryKeys.pagination(1),
      });
    },
    onSuccess: (data) => {
      toast(data, {
        icon: '👏',
        type: 'success',
      });
    },
    onError: (err: any) => {
      toast(err.detail, { icon: '😢', type: 'error' });
    },
    onSettled: async () => {},
  });
}
