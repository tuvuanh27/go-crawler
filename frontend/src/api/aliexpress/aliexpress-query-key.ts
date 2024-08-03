export const aliexpressQueryKeys = {
  all: ['aliexpress'],
  details: () => [...aliexpressQueryKeys.all, 'one'],
  detail: (id: string) => [...aliexpressQueryKeys.details(), id],
  exportExcel: () => [...aliexpressQueryKeys.all, 'export-excel'],
  pagination: (page: number) => [
    ...aliexpressQueryKeys.all,
    'pagination',
    page,
  ],
};
