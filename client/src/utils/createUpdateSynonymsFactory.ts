import { Synonym, Word } from "../types/general";

// This function will simply iterate over values from a modal and turn them into a Synonym object.
// Whatever exists will have an id with a number type which will be treated as a "Update" (since the record already exists by virtue of having a primary key).
// Whatever doesn't exist will have a null id and will be created as a new record in the DB.

export function createUpdateSynonymsFactory(
  existingSynonyms: Synonym[],
  inputSynonyms: string[],
  word_id: number,
) {
  let synonymsToKeep: Synonym[] = [];

  inputSynonyms.forEach((syn) => {
    for (const existingSynonym of existingSynonyms) {
      if (syn === existingSynonym.synonym) {
        synonymsToKeep.push({
          synonym: existingSynonym.synonym,
          id: existingSynonym.id,
          word_id,
          word: null
        });
        break;
      }
    }
    synonymsToKeep.push({ synonym: syn, id: null, word_id, word: null });
  });
  return synonymsToKeep;
}
