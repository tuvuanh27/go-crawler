import { useMutation } from '@tanstack/react-query';
import { apiClient } from '../api-client.ts';
import { ProductTypeSource } from '../../types/product.ts';
import { convertToUnixTime } from '../../helper';
import { toast } from 'react-toastify';

type Props = {
  productType: ProductTypeSource;
  startDate?: string | null;
  endDate?: string | null;
  exportType: 'full' | 'main';
};

export function useExportAliexpressProducts({
  productType,
  startDate,
  endDate,
  exportType,
}: Props) {
  const getExportProductsFn = async () => {
    let url = `product/export-excel-aliexpress?export_type=${exportType}&type=${productType}`;

    if (startDate && endDate) {
      url += `&start_date=${convertToUnixTime(
        startDate,
      )}&end_date=${convertToUnixTime(endDate)}`;
    }

    const response = await apiClient.get(url, {
      responseType: 'blob',
    });
    return response.data;
  };

  return useMutation({
    mutationFn: getExportProductsFn,
    onSuccess: (blob) => {
      // Create a URL for the Blob
      const url = window.URL.createObjectURL(
        new Blob([blob], {
          type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        }),
      );

      // Create a link element
      const link = document.createElement('a');

      // Set the download URL and file name
      link.href = url;
      link.setAttribute('download', 'exported-products.xlsx');

      // Append link to body and click to start download
      document.body.appendChild(link);
      link.click();

      // Clean up and remove the link
      link.parentNode?.removeChild(link);

      toast('Exported successfully', {
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
