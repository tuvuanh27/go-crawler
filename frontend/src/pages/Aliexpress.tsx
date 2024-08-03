import Breadcrumb from '../components/Breadcrumbs/Breadcrumb';
import React, { useEffect, useState } from 'react';
import { ProductTypeSource } from '../types/product.ts';
import {
  useCrawlAliProducts,
  usePaginatedAliexpressProducts,
} from '../api/aliexpress';
import ProductTable from '../components/Tables/ProductTable.tsx';
import { useExportAliexpressProducts } from '../api/aliexpress/use-export-aliexpress-products.ts';

const Aliexpress = () => {
  const [page, setPage] = useState(1);
  const [productIds, setProductIds] = useState<string>('');
  const [startDate, setStartDate] = useState<string | null>(null);
  const [endDate, setEndDate] = useState<string | null>(null);
  const [exportType, setExportType] = useState<'full' | 'main'>('full');
  const pageLimit = 20;

  const {
    data: paginatedProducts,
    refetch,
    error,
    isLoading,
  } = usePaginatedAliexpressProducts({
    page,
    pageLimit,
    productType: ProductTypeSource.Aliexpress,
    startDate,
    endDate,
  });

  const { mutate } = useExportAliexpressProducts({
    productType: ProductTypeSource.Aliexpress,
    startDate,
    endDate,
    exportType,
  });

  useEffect(() => {
    refetch();
  }, [startDate, endDate, refetch]);

  const crawlAliProducts = useCrawlAliProducts();

  // Ensure total is defined and handle possible undefined or null values
  const totalItems = paginatedProducts?.total || 0;
  const totalPages = Math.ceil(totalItems / pageLimit);

  const prevPage = () => {
    if (page > 1) setPage(page - 1);
  };

  const nextPage = () => {
    setPage(page + 1);
  };

  const handlePageClick = (pageNumber: number) => {
    setPage(pageNumber);
  };

  const getPageNumbers = () => {
    if (totalPages <= 6) {
      return Array.from({ length: totalPages }, (_, index) => index + 1);
    } else {
      let pages = [];
      if (page <= 4) {
        pages = [1, 2, 3, 4, 5, totalPages];
      } else if (page >= totalPages - 3) {
        pages = [
          1,
          totalPages - 4,
          totalPages - 3,
          totalPages - 2,
          totalPages - 1,
          totalPages,
        ];
      } else {
        pages = [1, page - 1, page, page + 1, totalPages];
      }
      return Array.from(new Set(pages)); // Remove duplicates
    }
  };

  const handleSetProductIds = (value: React.SetStateAction<string>) => {
    setProductIds(value);
  };

  const handleCrawlProducts = () => {
    if (productIds) {
      const ids = productIds
        .trim()
        .split(' ')
        .filter((id) => id !== '');
      crawlAliProducts.mutate(ids);
      setProductIds('');
    }
  };

  const changeStartDate = (e: {
    target: { value: React.SetStateAction<string | null> };
  }) => {
    setStartDate(e.target.value);
  };

  const changeEndDate = (e: {
    target: { value: React.SetStateAction<string | null> };
  }) => {
    setEndDate(e.target.value);
  };

  const onChangeExportType = (e: any) => {
    setExportType(e.target.value);
  };

  const handleExportExcel = () => {
    mutate();
  };

  return (
    <>
      <Breadcrumb
        pageName={`Aliexpress: ${paginatedProducts?.total || 0} Products`}
      />

      <div className="grid grid-cols-1 gap-2 sm:grid-cols-12">
        <div className="flex flex-col gap-2 col-span-4">
          <div className="flex flex-col gap-5.5 p-6.5 ">
            <div className="mb-5">
              <label className="mb-3 block text-black dark:text-white">
                Product Ids
              </label>
              <div className="relative flex items-center">
                <input
                  type="text"
                  placeholder="Product Ids"
                  className="w-full rounded-lg border-[1.5px] border-stroke bg-transparent py-3 px-5 pr-20 text-black outline-none transition focus:border-primary active:border-primary disabled:cursor-default disabled:bg-whiter dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
                  value={productIds}
                  onChange={(e) => handleSetProductIds(e.target.value)}
                />
                <button
                  onClick={handleCrawlProducts}
                  className="absolute right-2 top-1/2 transform -translate-y-1/2 bg-blue-500 text-white rounded-lg px-4 py-2 hover:bg-blue-600 focus:outline-none"
                >
                  Add
                </button>
              </div>
            </div>
          </div>
        </div>
        <div className="flex flex-col gap-2 col-span-2">
          <div className="flex flex-col gap-5.5 p-6.5">
            <div>
              <label className="mb-3 block text-black dark:text-white">
                From Date
              </label>
              <input
                type="date"
                onChange={changeStartDate}
                className="w-full rounded-lg border-[1.5px] border-stroke bg-transparent py-3 px-5 text-black outline-none transition focus:border-primary active:border-primary disabled:cursor-default disabled:bg-white dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
              />
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-2 col-span-2">
          <div className="flex flex-col gap-5.5 p-6.5">
            <div>
              <label className="mb-3 block text-black dark:text-white">
                To Date
              </label>
              <input
                type="date"
                onChange={changeEndDate}
                className="w-full rounded-lg border-[1.5px] border-stroke bg-transparent py-3 px-5 text-black outline-none transition focus:border-primary active:border-primary disabled:cursor-default disabled:bg-white dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
              />
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-2 col-span-2">
          <div className="flex flex-col gap-5.5 p-6.5">
            <div>
              <label className="mb-3 block text-black dark:text-white">
                Export Type
              </label>
              <select
                onChange={onChangeExportType}
                value={exportType}
                className="w-full rounded-lg border-[1.5px] border-stroke bg-transparent py-3 px-5 text-black outline-none transition focus:border-primary active:border-primary disabled:cursor-default disabled:bg-white dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
              >
                <option value="full">Full</option>
                <option value="main">Main</option>
              </select>
            </div>
          </div>
        </div>

        <div className={`flex flex-col gap-2 col-span-2`}>
          <div className="flex flex-col gap-5.5 p-6.5">
            <div>
              <label className="mb-3 block text-black dark:text-white">
                Export
              </label>
              <button
                className="flex w-full justify-center rounded bg-primary p-3 font-medium text-gray hover:bg-opacity-80"
                onClick={handleExportExcel}
              >
                <span className="pr-2">Download</span>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="size-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      {error instanceof Error && <div>{error.message}</div>}

      {isLoading && <div>Loading...</div>}

      <ProductTable products={paginatedProducts?.products || []} />

      <div>
        <div className="flex justify-center my-4">
          <nav className="inline-flex space-x-2" aria-label="Pagination">
            <button
              onClick={prevPage}
              disabled={page <= 1}
              className="hover:cursor-pointer inline-flex items-center justify-center w-10 h-10 bg-gray-100 text-gray-600 hover:bg-blue-100 hover:text-blue-600 rounded-full"
            >
              <span className="sr-only">Previous</span>
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M15 19l-7-7 7-7"
                ></path>
              </svg>
            </button>

            {getPageNumbers().map((number, index, array) => (
              <React.Fragment key={number}>
                {index > 0 && array[index - 1] + 1 < number && (
                  <span className="inline-flex items-center justify-center w-10 h-10">
                    ...
                  </span>
                )}
                <button
                  onClick={() => handlePageClick(number)}
                  className={`inline-flex items-center justify-center w-10 h-10 ${
                    number === page
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 text-gray-600 hover:bg-blue-100 hover:text-blue-600'
                  } rounded-full`}
                >
                  {number}
                </button>
              </React.Fragment>
            ))}

            <button
              onClick={nextPage}
              disabled={
                paginatedProducts?.products &&
                paginatedProducts?.products.length < pageLimit
              }
              className="hover:cursor-pointer inline-flex items-center justify-center w-10 h-10 bg-gray-100 text-gray-600 hover:bg-blue-100 hover:text-blue-600 rounded-full"
            >
              <span className="sr-only">Next</span>
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M9 5l7 7-7 7"
                ></path>
              </svg>
            </button>
          </nav>
        </div>
      </div>
    </>
  );
};

export default Aliexpress;
