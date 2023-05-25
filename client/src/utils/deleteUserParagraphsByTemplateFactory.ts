import { Paragraph } from "../types/general";

export default function deleteUserParagraphsByTemplateFactory(
  values: { input: string },
  userParagraphsByWord: Paragraph[]
): { delete_paragraphs: number[] } {
  let delete_paragraphs: number[] = [];
  const vals = values.input.split("\n");

  filter: for (let i = 0; i < userParagraphsByWord.length; i++) {
    for (let n = 0; n < vals.length; n++) {
      if (userParagraphsByWord[i].name === vals[n]) {
        continue filter;
      }
    }
    // Any existing synonym that is not in the newly formed array means it's going to be deleted.
    delete_paragraphs.push(userParagraphsByWord[i].id!);
  }

  return { delete_paragraphs };
}
