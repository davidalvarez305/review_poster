import { Sentence } from "../types/general";

export function createUpdateSentences(
  existingSentences: Sentence[],
  inputSentences: string[],
  template_id: number,
  paragraph_id: number,
  user_id: number
): Sentence[] {
  let sentencesToKeep: Sentence[] = [];

  inputSentences.forEach((sentence) => {
    for (const existingSentence of existingSentences) {
      if (sentence === existingSentence.sentence) {
        sentencesToKeep.push({
          sentence: existingSentence.sentence,
          id: existingSentence.id,
          template_id,
          paragraph_id,
          user_id
        });
        break;
      }
    }
    sentencesToKeep.push({ sentence, id: null, template_id, paragraph_id, user_id });
  });
  return sentencesToKeep;
}
