const findPattern = (pattern) => {
  const files = glob(pattern);
  if (files.length === 0) {
    return [];
  }
  return files;
}