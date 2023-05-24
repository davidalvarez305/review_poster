import { Word, WordFormInput } from "../types/general";
import createTagFactory from "./createTagFactory";
import { getId } from "./getId";

export default function createWordsFactory(
  values: WordFormInput,
  user_id: number,
  words: Word[]
): Word {
  return {
    id: getId(values.name, words, "name"),
    tag: createTagFactory(values.name),
    name: values.name,
    user_id,
    user: null,
    synonyms: null,
  };
}
