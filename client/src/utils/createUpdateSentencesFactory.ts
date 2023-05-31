import { Sentence } from "../types/general";

export default function createUpdateSentencesFactory(
  values: { input: string },
  userSentencesByParagraphAndTemplate: Sentence[]
): { sentences: Sentence[] } {
  let sentences: Sentence[] = [];
  const vals = values.input.split("\n");

  add: for (let n = 0; n < vals.length; n++) {
    for (let i = 0; i < userSentencesByParagraphAndTemplate.length; i++) {
      if (userSentencesByParagraphAndTemplate[i].sentence === vals[n]) {
        sentences.push({
          ...userSentencesByParagraphAndTemplate[i],
          paragraph: null,
        });
        continue add;
      }
    }

    if (vals[n].length > 0) {
      sentences.push({
        id: null,
        sentence: vals[n],
        paragraph_id: userSentencesByParagraphAndTemplate[0].paragraph_id,
        paragraph: null,
      });
    }
  }

  return { sentences };
}
