import { Synonym } from "../types/general";

export default function deleteUserSynonymsByWordFactory(
  values: { input: string },
  userSynonymsByWord: Synonym[]
): { delete_synonyms: number[] } {
  let delete_synonyms: number[] = [];
  const vals = values.input.split("\n");

  filter: for (let i = 0; i < userSynonymsByWord.length; i++) {
    for (let n = 0; n < vals.length; n++) {
      if (userSynonymsByWord[i].synonym === vals[n]) {
        continue filter;
      }
    }
    // Any existing synonym that is not in the newly formed array means it's going to be deleted.
    delete_synonyms.push(userSynonymsByWord[i].id!);
  }

  return { delete_synonyms };
}
