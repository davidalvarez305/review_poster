import { Paragraph, Template } from "../types/general";

export function createUpdateParagraphsFactory(
  existingParagraphs: Paragraph[],
  inputParagraphs: string[],
  template_id: number,
  template: Template,
): Paragraph[] {
  let paragraphsToKeep: Paragraph[] = [];

  inputParagraphs.forEach((paragraph) => {
    for (const existingParagraph of existingParagraphs) {
      if (paragraph === existingParagraph.name) {
        paragraphsToKeep.push({
          name: existingParagraph.name,
          id: existingParagraph.id,
          template_id,
          order: existingParagraph.order,
          template,
        });
        break;
      }
    }
    paragraphsToKeep.push({
      name: paragraph,
      id: null,
      template_id,
      template,
    });
  });
  return paragraphsToKeep;
}