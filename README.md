# protomask

protomask is a package that lets you update protobuf messages with a help of field masks.

**Features**

- Assign values from one message to another based on a field mask.
- Create a mask of all populated fields of a message. 
- Supports nested fields, and handles their parent being nil.
- Supports [fieldmaskpb](https://google.golang.org/protobuf/types/known/), including `oneof` properties.

## Installation

`go get github.com/olexnzarov/protomask`

## Usage

```go
func (s *BookServer) UpdateBookPrice(id int64, priceCents int64) error {
  update := &pb.Book{
    Id: id,
    Price: &pb.Price{
      Cents: priceCents,
    },
  }
  mask, _ := fieldmaskpb.New(update, "price.cents")
  return s.UpdateBook(update, mask)
}

func (s *BookServer) UpdateBook(update *pb.Book, updateMask protomask.FieldMask) error {
  book, err := s.bookStorage.GetById(update.Id)
  if err != nil {
    return err
  }
  err = protomask.Update(book, update, updateMask)
  if err != nil {
    return err
  }
  return s.bookStorage.Save(book)
}
```

## Examples

See [protomask_test.go](./protomask_test.go) for more examples on how to use the package. Also, you can check out [gofu](https://github.com/olexnzarov/gofu) for a real-life example of how it can be used.