import { Paragraph, Sentence } from "../types/general";

export function createUpdateSentencesFactory(
  existingSentences: Sentence[],
  inputSentences: string[],
  paragraph_id: number,
  paragraph: Paragraph
): Sentence[] {
  let sentencesToKeep: Sentence[] = [];

  inputSentences.forEach((sentence) => {
    for (const existingSentence of existingSentences) {
      if (sentence === existingSentence.sentence) {
        sentencesToKeep.push({
          sentence: existingSentence.sentence,
          id: existingSentence.id,
          paragraph_id,
          paragraph,
        });
        break;
      }
    }
    sentencesToKeep.push({
      sentence,
      id: null,
      paragraph_id,
      paragraph,
    });
  });
  return sentencesToKeep;
}
