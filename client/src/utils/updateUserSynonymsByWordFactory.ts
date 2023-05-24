import { Synonym, UpdateUserSynonymsByWordInput } from "../types/general";

export default function updateUserSynonymsByWordFactory(
  values: { input: string },
  userSynonymsByWord: Synonym[]
): UpdateUserSynonymsByWordInput {
  let delete_synonyms: number[] = [];
  let synonyms: Synonym[] = [];

  const vals = values.input.split("\n");

  for (let n = 0; n < vals.length; n++) {
    for (let i = 0; i < userSynonymsByWord.length; i++) {
      if (userSynonymsByWord[i].synonym === vals[n]) {
        synonyms.push(userSynonymsByWord[i]);
      } else {
        synonyms.push({
          id: null,
          synonym: vals[n],
          word_id: userSynonymsByWord[0].word_id,
          word: null,
        });
      }
    }
  }

  filter:
  for (let i = 0; i < userSynonymsByWord.length; i++) {
    for (let n = 0; n < synonyms.length; n++) {
      if (userSynonymsByWord[i].synonym === synonyms[n].synonym) {
        continue filter;
      }
    }
    // Any existing synonym that is not in the newly formed array means it's going to be deleted.
    delete_synonyms.push(userSynonymsByWord[i].id!);
  }

  return { synonyms, delete_synonyms };
}
