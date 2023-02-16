export function extractTags(tag: string) {
  const regex = /[A-Za-z0-9]+/g;
  const match = tag.match(regex);
  if (match) {
    return { word: match[0], tag: "(#" + match[0] + ")" };
  }
  return { word: "", tag: "" };
}
