import { Synonym } from "../types/general";

export default function updateUserSynonymsByWordFactory(
  values: { input: string },
  userSynonymsByWord: Synonym[]
): { synonyms: Synonym[] } {
  let synonyms: Synonym[] = [];
  const vals = values.input.split("\n");

  add: for (let n = 0; n < vals.length; n++) {
    for (let i = 0; i < userSynonymsByWord.length; i++) {
      if (userSynonymsByWord[i].synonym === vals[n]) {
        synonyms.push({
          ...userSynonymsByWord[i],
          word: null,
        });
        continue add;
      }
    }

    if (vals[n].length > 0) {
      synonyms.push({
        id: null,
        synonym: vals[n],
        word_id: userSynonymsByWord[0].word_id,
        word: null,
      });
    }
  }

  return { synonyms };
}
