import { Sentence } from "../types/general";

export default function deleteUserParagraphSentencesByTemplateFactory(
  values: { input: string },
  userParagraphSentencesByTemplate: Sentence[]
): { delete_sentences: number[] } {
  let delete_sentences: number[] = [];
  const vals = values.input.split("\n");

  filter: for (let i = 0; i < userParagraphSentencesByTemplate.length; i++) {
    for (let n = 0; n < vals.length; n++) {
      if (userParagraphSentencesByTemplate[i].sentence === vals[n]) {
        continue filter;
      }
    }
    // Any existing synonym that is not in the newly formed array means it's going to be deleted.
    delete_sentences.push(userParagraphSentencesByTemplate[i].id!);
  }

  return { delete_sentences };
}
