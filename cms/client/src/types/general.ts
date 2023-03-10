export type Token = {
  id: number;
  uuid: string;
  created_at: number;
};

export type User = {
  id: number;
  username: string;
  password: string;
  email: string;
  token: Token | null;
};

export type Template = {
  id: number;
  name: string;
  user_id: number;
  user: User | null;
};

export type Word = {
  id: number;
  name: string;
  tag: string;
  user_id: number;
  user: User | null;
  synonyms: Synonym[] | null;
};

export type Synonym = {
  id: number | null;
  synonym: string;
  word_id: number;
  word: Word | null;
};

export type Sentence = {
  id: number | null;
  sentence: string;
  paragraph_id: number;
  template_id: number;
  user_id: number;
  paragraph: Paragraph | null;
  template: Template | null;
};

export type Paragraph = {
  id: number | null;
  name: string;
  order?: number;
  template_id: number;
  user_id: number;
  template: Template | null;
  user: User | null;
};

export type UpdateParagraph = {
  id?: number;
  name: string;
  order: number | null;
  template_id: number;
};

export type UpdateSynonym = {
  id?: number;
  synonym: string;
  word_id: number;
};

export type UpdateSentence = {
  id?: number;
  sentence: string;
  paragraph_id: number;
  template_id: number;
};

export type Content = {
  sentences: string;
  template: string;
  paragraph: string;
  order: number;
};

export type FinalizedContent = {
  sentences: string[];
  template: string;
  paragraph: string;
  order: number;
};

export type DictionaryResponse = {
  name: string;
  tag: string;
  synonyms: string;
};

export type Dictionary = {
  [key: string]: string[];
};
