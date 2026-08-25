# Services

## Standards & Practices

### Invert if-statements

#### Bad example

```go
if condition1 {
    if condition2 {
        return // something
    } else {
        return // something
    }
}
```

#### Good example

```go
if !condition1 {
    return // something

}
if condition2 {
    return // something
}
return // something
```

### Separation of concerns

#### Bad example

This is all in the repository. It has logic and also inlines the queries

**repository.go**

```go
func upsertEntry(
 ctx context.Context,
 tx pgx.Tx,
 accountID int32,
 listID int32,
 entry domain.TranslationEntryInput,
) (int32, error) {
 if entry.ID != nil {
  tag, err := tx.Exec(ctx, `
   UPDATE translation
   SET list_id = $1, original_html = $2, translation_html = $3, created_at = $4, updated_at = $5
   WHERE id = $6 AND list_id IN (SELECT id FROM translation_list WHERE account_id = $7)
  `, listID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt, entry.UpdatedAt, *entry.ID, accountID)
  if err != nil {
   return 0, err
  }
  if tag.RowsAffected() == 1 {
   return *entry.ID, nil
  }
  return 0, domain.ErrEntryNotFound
 }

 var newID int32
 err := tx.QueryRow(ctx, `
  INSERT INTO translation (list_id, original_html, translation_html, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING id
 `, listID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt, entry.UpdatedAt).Scan(&newID)
 return newID, err
}
```

#### Good example

**user_case.go**

Handles business logic etc.

```go
func upsertEntry(
 ctx context.Context,
 tx pgx.Tx,
 accountID int32,
 listID int32,
 entry domain.TranslationEntryInput,
) (int32, error) {
 if entry.ID != nil {
  rowsAffected, err := repository.UpdateTranslationEntry(ctx, tx, accountID, listID, *entry.ID, entry)
  if err != nil {
   return 0, err
  }
  if rowsAffected == 1 {
   return *entry.ID, nil
  }
  return 0, domain.ErrEntryNotFound
 }

 return repository.InsertTranslationEntry(ctx, tx, listID, entry)
}
```

**repository.go**

Only responsible for database actions

```go
func UpdateTranslationEntry(
 ctx context.Context,
 tx pgx.Tx,
 accountID int32,
 listID int32,
 entryID int32,
 entry domain.TranslationEntryInput,
) (int64, error) {
 tag, err := tx.Exec(ctx, `
  UPDATE translation
  SET list_id = $1, original_html = $2, translation_html = $3, created_at = $4, updated_at = $5
  WHERE id = $6 AND list_id IN (SELECT id FROM translation_list WHERE account_id = $7)
 `, listID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt, entry.UpdatedAt, entryID, accountID)
 if err != nil {
  return 0, err
 }
 return tag.RowsAffected(), nil
}

func InsertTranslationEntry(
 ctx context.Context,
 tx pgx.Tx,
 listID int32,
 entry domain.TranslationEntryInput,
) (int32, error) {
 var newID int32
 err := tx.QueryRow(ctx, `
  INSERT INTO translation (list_id, original_html, translation_html, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING id
 `, listID, entry.OriginalHTML, entry.TranslationHTML, entry.CreatedAt, entry.UpdatedAt).Scan(&newID)
 return newID, err
}
```
