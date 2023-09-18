package protomask

import (
	"testing"
	"time"

	"github.com/olexnzarov/protomask/internal/pbtest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUpdate(t *testing.T) {
	book := &pbtest.Book{
		Id:   1605,
		Name: "Don Quixote",
		Price: &pbtest.Price{
			Cents: 1500,
		},
	}
	bookReference := proto.Clone(book).(*pbtest.Book)

	bookUpdate := &pbtest.Book{
		Name: "Don Quixote: Special Edition",
		Price: &pbtest.Price{
			Cents: 1000,
			Discount: &pbtest.Discount{
				FullPrice: &pbtest.Price{
					Cents: book.Price.Cents,
				},
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		},
	}
	mask, _ := fieldmaskpb.New(bookUpdate, "price", "name")
	err := Update(book, bookUpdate, mask)
	if err != nil {
		t.Fatalf("failed to update: %s", err)
		return
	}

	if book.Price.Cents != bookUpdate.Price.Cents ||
		book.Price.Discount.FullPrice != bookUpdate.Price.Discount.FullPrice ||
		book.Price.Discount.ExpiresAt != bookUpdate.Price.Discount.ExpiresAt ||
		book.Name != bookUpdate.Name {
		t.Fatal("assertion after the update failed: message was not updated")
		return
	}

	if bookReference.Id != book.Id {
		t.Fatal("assertion after the update failed: properties outside the mask were updated")
		return
	}
}

func TestUpdateSetNilValue(t *testing.T) {
	book := &pbtest.Book{
		Price: &pbtest.Price{
			Cents: 1000,
			Discount: &pbtest.Discount{
				FullPrice: &pbtest.Price{
					Cents: 1500,
				},
				ExpiresAt: time.Now().Unix(),
			},
		},
	}

	bookUpdate := &pbtest.Book{}
	mask, _ := fieldmaskpb.New(bookUpdate, "price.discount")
	err := Update(book, bookUpdate, mask)
	if err != nil {
		t.Fatalf("failed to update: %s", err)
		return
	}

	if book.Price.Discount != nil {
		t.Fatal("assertion after the update failed: message was not updated")
		return
	}
}

func TestUpdateSetValueToNilParent(t *testing.T) {
	book := &pbtest.Book{}

	bookUpdate := &pbtest.Book{
		Price: &pbtest.Price{
			Discount: &pbtest.Discount{
				FullPrice: &pbtest.Price{
					Cents: 1500,
				},
			},
		},
	}
	mask, _ := fieldmaskpb.New(bookUpdate, "price.discount.full_price.cents")
	err := Update(book, bookUpdate, mask)
	if err != nil {
		t.Fatalf("failed to update: %s", err)
		return
	}

	if book.Price.Discount.FullPrice.Cents != bookUpdate.Price.Discount.FullPrice.Cents {
		t.Fatal("assertion after the update failed: message was not updated")
		return
	}
}

func TestOneOf(t *testing.T) {
	priceReply := &pbtest.PriceReply{
		Response: &pbtest.PriceReply_Error{
			Error: &pbtest.Error{Message: "unknown price"},
		},
	}

	priceReplyUpdate := &pbtest.PriceReply{
		Response: &pbtest.PriceReply_Price{
			Price: &pbtest.Price{
				Cents: 1500,
			},
		},
	}
	mask, err := fieldmaskpb.New(priceReplyUpdate, "price")
	if err != nil {
		t.Fatalf("failed to create a field mask: %s", err)
		return
	}
	err = Update(priceReply, priceReplyUpdate, mask)
	if err != nil {
		t.Fatalf("failed to update: %s", err)
		return
	}

	if priceReply.GetPrice() == nil || priceReply.GetPrice().Cents != priceReplyUpdate.GetPrice().Cents {
		t.Fatal("assertion after the update failed: message was not updated")
		return
	}
}

type invalidFieldMask struct {
	paths []string
}

func (*invalidFieldMask) IsValid(m protoreflect.ProtoMessage) bool {
	return true
}

func (mask *invalidFieldMask) GetPaths() []string {
	return mask.paths
}

func TestInvalidFieldMask(t *testing.T) {
	book := &pbtest.Book{}

	bookUpdate := &pbtest.Book{
		Name: "Null References: The Billion Dollar Mistake",
	}
	mask := &invalidFieldMask{paths: []string{"Name"}}
	err := Update(book, bookUpdate, mask)
	if err == nil {
		t.Fatal("Update did not return error on invalid field mask")
		return
	}
}
