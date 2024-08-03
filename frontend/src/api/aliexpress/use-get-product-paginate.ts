import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api-client.ts';
import { aliexpressQueryKeys } from './aliexpress-query-key.ts';
import { Product, ProductTypeSource } from '../../types/product.ts';
import { convertToUnixTime } from '../../helper';

type Props = {
  page: number;
  pageLimit: number;
  productType: ProductTypeSource;
  startDate?: string | null;
  endDate?: string | null;
};

type PaginatedAliexpressProductsResponse = {
  products: Product[];
  total: number;
  limit: number;
  page: number;
};

export function usePaginatedAliexpressProducts({
  page,
  pageLimit,
  productType,
  startDate,
  endDate,
}: Props) {
  const getPaginatedProductsFn =
    async (): Promise<PaginatedAliexpressProductsResponse> => {
      let url = `product/get-products?page=${page}&limit=${pageLimit}&type=${productType}`;

      if (startDate && endDate) {
        url += `&start_date=${convertToUnixTime(
          startDate,
        )}&end_date=${convertToUnixTime(endDate)}`;
      }

      const response = await apiClient.get(url);
      return response.data;
    };

  return useQuery({
    queryKey: aliexpressQueryKeys.pagination(page),
    queryFn: () => getPaginatedProductsFn(),
  });
}
