// Minimal shims to unblock TS complaining about Next internal modules in generated .next/types
declare module 'next/dist/lib/metadata/types/metadata-interface.js' {
  export type ResolvingMetadata = any;
  export type ResolvingViewport = any;
}

declare module 'next/navigation' {
  export const useRouter: any;
  export const useParams: any;
  export const useSearchParams: any;
  export const redirect: any;
}

declare module 'next/link' {
  const Link: any;
  export default Link;
}

declare module 'next/image' {
  const Image: any;
  export default Image;
}

declare module 'next/font/google' {
  export const Inter: any;
  export const Roboto: any;
}
