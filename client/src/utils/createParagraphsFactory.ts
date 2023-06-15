import { Paragraph } from "../types/general";

export default function createParagraphsFactory(values: {
  paragraphs: string;
  template_id: number;
}): Paragraph[] {
  return values.paragraphs.split("\n").map((i) => {
    return {
      name: i,
      template_id: values.template_id,
      id: null,
      template: null,
      sentences: null,
    };
  });
}
