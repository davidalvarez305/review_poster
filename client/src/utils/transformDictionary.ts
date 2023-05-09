import { Dictionary, DictionaryResponse } from "../types/general";

export function transformDictionary(dict: DictionaryResponse[]): Dictionary {
  let d: Dictionary = {};
  dict.forEach((ea) => {
    d[ea.tag] = ea.synonyms.split("///");
  });
  return d;
}
