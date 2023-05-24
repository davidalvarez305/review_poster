import { Synonym, Word } from "../types/general";

export default function createSynonymsFactory(
  values: { input: string },
  word: Word
): Synonym[] {
  const vals = values.input.split("\n");
  let synonyms: Synonym[] = [];

  for (let i = 0; i < vals.length; i++) {
    synonyms.push({
      id: null,
      word_id: word.id!,
      synonym: vals[i],
      word: null,
    });
  }

  return synonyms;
}
