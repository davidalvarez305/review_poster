import { Paragraph } from "../types/general";

export default function updateUserParagraphsByTemplateFactory(
  values: { input: string },
  userParagraphsByTemplate: Paragraph[]
): { paragraphs: Paragraph[] } {
  let paragraphs: Paragraph[] = [];
  const vals = values.input.split("\n");

  add: for (let n = 0; n < vals.length; n++) {
    for (let i = 0; i < userParagraphsByTemplate.length; i++) {
      if (userParagraphsByTemplate[i].name === vals[n]) {
        paragraphs.push({
          ...userParagraphsByTemplate[i],
          template: null,
        });
        continue add;
      }
    }

    if (vals[n].length > 0) {
      paragraphs.push({
        id: null,
        name: vals[n],
        template_id: userParagraphsByTemplate[0].template_id,
        template: null,
      });
    }
  }

  return { paragraphs };
}
