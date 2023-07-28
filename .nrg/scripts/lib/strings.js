const sortstrings = (strings) => {
  const sorted = strings.sort((a, b) => {
    const aLC = a.toLowerCase();
    const bLC = b.toLowerCase();
    return aLC.localeCompare(bLC);
  });
  return sorted;
}

const longeststring = (strings) => {
    let longest = 0;
    strings.forEach((element) => {
        if (element.length > longest) {
            longest = element.length;
        }
    });
    return longest;
}