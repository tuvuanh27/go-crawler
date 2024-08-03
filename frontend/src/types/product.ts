export enum ProductTypeSource {
  Aliexpress = 1,
  Amazon = 2,
  Ebay = 3,
}

export const ValidProductTypeSources: Record<number, ProductTypeSource> = {
  1: ProductTypeSource.Aliexpress,
  2: ProductTypeSource.Amazon,
  3: ProductTypeSource.Ebay,
};

export interface Image {
  url: string;
  z_index: number; // The order of the image, the first image with z_index = 0 is the main image
}

export interface VariationType {
  valueId: string;
  skuPropId: string;
  name: string;
  image: string;
}

export interface Variation {
  sizes: VariationType[];
  colors: VariationType[];
}

export interface Specification {
  name: string;
  value: string;
}

export interface Sku {
  skuId: string;
  skuAttr: string;
  price: string;
  promotionPrice: string;
}

export interface Seller {
  storeId: string;
  storeName: string;
  shippingRating: string;
  communicationRating: string;
  itemAsDescribed: string;
}

export interface Product {
  _id?: string; // Optional field as in Go `omitempty`
  productId: string;
  title: string;
  description: string;
  specifications: Specification[];
  productTypeSource: ProductTypeSource;
  skus: Sku[];
  images: Image[];
  price: number;
  originalPrice: number;
  variation: Variation;
  seller: Seller;
  createdAt: Date;
  updatedAt: Date;
}
