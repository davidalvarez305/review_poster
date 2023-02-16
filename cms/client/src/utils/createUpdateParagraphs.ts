import { Paragraph } from "../types/general";

export function createUpdateParagraphs(
  existingParagraphs: Paragraph[],
  inputParagraphs: string[],
  template_id: number,
  user_id: number
): Paragraph[] {
  let paragraphsToKeep: Paragraph[] = [];

  inputParagraphs.forEach((paragraph) => {
    for (const existingParagraph of existingParagraphs) {
      if (paragraph === existingParagraph.name) {
        paragraphsToKeep.push({
          name: existingParagraph.name,
          id: existingParagraph.id,
          template_id,
          user_id,
          order: existingParagraph.order
        });
        break;
      }
    }
    paragraphsToKeep.push({ name: paragraph, id: null, template_id, user_id });
  });
  return paragraphsToKeep;
}
