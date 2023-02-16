import { Content, FinalizedContent } from "../types/general";

export function transformContent(content: Content[]): FinalizedContent[] {
  return content.map((c) => {
    return { ...c, sentences: c.sentences.split("///") };
  });
}
