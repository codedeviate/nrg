const findPaths = (pathname) => {
  const files = finddirname(pathname);
  if (files.length === 0) {
    return [];
  }
  return files;
}