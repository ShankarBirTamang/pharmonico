// Package database provides index management for MongoDB collections
package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates all necessary indexes for MongoDB collections
func (mc *MongoClient) CreateIndexes(ctx context.Context) error {
	if err := mc.createPharmacyIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create pharmacy indexes: %w", err)
	}

	if err := mc.createPrescriberIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create prescriber indexes: %w", err)
	}

	if err := mc.createPatientIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create patient indexes: %w", err)
	}

	return nil
}

// createPharmacyIndexes creates indexes for the pharmacies collection
func (mc *MongoClient) createPharmacyIndexes(ctx context.Context) error {
	collection := mc.GetCollection("pharmacies")

	indexes := []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"ncpdp_id": 1},
			Options: options.Index().SetUnique(true).SetName("idx_ncpdp_id"),
		},
		{
			Keys:    map[string]interface{}{"name": 1},
			Options: options.Index().SetName("idx_name"),
		},
		{
			Keys:    map[string]interface{}{"address.zip": 1},
			Options: options.Index().SetName("idx_zip"),
		},
		{
			Keys:    map[string]interface{}{"location": "2dsphere"},
			Options: options.Index().SetName("idx_location_2dsphere"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createPrescriberIndexes creates indexes for the prescribers collection
func (mc *MongoClient) createPrescriberIndexes(ctx context.Context) error {
	collection := mc.GetCollection("prescribers")

	indexes := []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"npi": 1},
			Options: options.Index().SetUnique(true).SetName("idx_npi"),
		},
		{
			Keys:    map[string]interface{}{"name": 1},
			Options: options.Index().SetName("idx_name"),
		},
		{
			Keys:    map[string]interface{}{"license_number": 1},
			Options: options.Index().SetName("idx_license_number"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createPatientIndexes creates indexes for the patients collection
func (mc *MongoClient) createPatientIndexes(ctx context.Context) error {
	collection := mc.GetCollection("patients")

	indexes := []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true).SetName("idx_email"),
		},
		{
			Keys:    map[string]interface{}{"phone": 1},
			Options: options.Index().SetName("idx_phone"),
		},
		{
			Keys:    map[string]interface{}{"date_of_birth": 1},
			Options: options.Index().SetName("idx_date_of_birth"),
		},
		{
			Keys:    map[string]interface{}{"created_at": 1},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
