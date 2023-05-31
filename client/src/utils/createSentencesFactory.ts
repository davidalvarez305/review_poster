import { Sentence, Paragraph } from "../types/general";

export default function createSentencesFactory(
  values: { input: string },
  paragraph: Paragraph
): Sentence[] {
  const vals = values.input.split("\n");
  let sentences: Sentence[] = [];

  for (let i = 0; i < vals.length; i++) {
    sentences.push({
      id: null,
      paragraph_id: paragraph.id!,
      sentence: vals[i],
      paragraph: null,
    });
  }

  return sentences;
}
