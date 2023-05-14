import { Dictionary, Word } from "../types/general";

export function createDictionaryFactory(dict: Word[]): Dictionary {
  let d: Dictionary = {};

  for (let i = 0; i < dict.length; i++) {
    const word = dict[i];
    if (word.synonyms) d[word.tag] = [...word.synonyms.map(synonym => synonym.synonym)];
  }

  return d;
}
