import { Sentence } from "../types/general";
import { re } from "./spinContent";

export const filterSentences = (sentences: Sentence[]) => {
  let filtered: Sentence[] = [];
  sentences.forEach((s) => {
    const split = s.sentence.match(/\b\w{4,}/g);
    const matches = s.sentence.match(re);
    if (split!.length > matches!.length + 20) {
      filtered.push(s);
    }
  });
  return filtered;
};
