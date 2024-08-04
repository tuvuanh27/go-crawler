import { Product, VariationType } from '../../types/product.ts';

type Props = {
  products: Product[];
};

export default function ProductTable({ products }: Props) {
  function getVariation(variation: VariationType[]): string {
    return (variation || [])
      .map((c) => c.name)
      .sort(
        (a, b) =>
          (a || '').localeCompare(b || '', undefined, { numeric: true }) || 0,
      )
      .join(', ');
  }

  function getPrice(product: Product): number {
    return (
      Math.round( product.price * 10) / 10 ||
      product.originalPrice ||
      +product.skus[0]?.promotionPrice ||
      0
    );
  }

  return (
    <>
      <div className="grid grid-cols-1 gap-9 sm:grid-cols-12">
        <div className="flex flex-col col-span-12 gap-5.5 p-6.5 overflow-x-auto">
          <div className="">
            <table className="min-w-full divide-y-2 divide-gray-200 text-sm">
              <thead className="ltr:text-left rtl:text-right text-left">
                <tr>
                  <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">
                    <span className="invisible">Placeholder</span>
                  </th>
                  <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">
                    Product Id
                  </th>
                  <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">
                    Title
                  </th>
                  <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">
                    Price
                  </th>
                  <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">
                    Size
                  </th>
                  <th className="whitespace-nowrap px-4 py-2 font-medium text-gray-900">
                    Color
                  </th>
                  <th className="px-4 py-2"></th>
                </tr>
              </thead>

              <tbody className="divide-y divide-gray-200">
                {products.map((product, key) => (
                  <tr key={key}>
                    <td className="px-4 py-2">
                      <img
                        src={product.images[0].url}
                        alt="Product"
                        className="size-16 rounded-lg object-cover shadow-sm"
                      />
                    </td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700">
                      <a
                        href={`https://www.aliexpress.us/item/${product.productId}.html`}
                        className="text-blue-500 hover:underline"
                        target="_blank"
                      >
                        {product.productId}
                      </a>
                    </td>

                    <td className="whitespace-nowrap px-4 py-2 text-gray-700">
                      {product.title.slice(0, 50) + '...'}
                    </td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700">
                      {getPrice(product)}
                    </td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700">
                      {getVariation(product.variation.sizes)}
                    </td>
                    <td className="whitespace-nowrap px-4 py-2 text-gray-700">
                      {getVariation(product.variation.colors)}
                    </td>
                    <td className="px-4 py-2">
                      <a
                        href="#"
                        className="inline-block rounded bg-indigo-600 px-4 py-2 text-xs font-medium text-white hover:bg-indigo-700"
                      >
                        View
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </>
  );
}
