// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"log"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// package validations
func ValidatePackage(pkg packages.Package) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidatePackageName(pkg.ObjectMeta.Name, field.NewPath("metadata").Child("name"))...)

	return allErrs
}

// validate name
func ValidatePackageName(pkgName string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs,
		validation.IsFullyQualifiedName(fldPath, pkgName)...)

	return allErrs
}

// package version validations

func ValidatePackageVersion(pv packages.PackageVersion) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs,
		ValidatePackageVersionSpecPackageName(pv.Spec.PackageName, field.NewPath("spec", "packageName"))...)

	allErrs = append(allErrs, ValidatePackageVersionSpecVersion(pv.Spec.Version, field.NewPath("spec", "version"))...)

	allErrs = append(allErrs,
		ValidatePackageVersionName(pv.ObjectMeta.Name, pv.Spec.PackageName, field.NewPath("metadata").Child("name"))...)

	return allErrs
}

// validate metdata.name = spec.PackageName + spec.Version
func ValidatePackageVersionName(pvName, pkgName string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !strings.HasPrefix(pvName, pkgName+".") {
		log.Printf("PVN: %s, PN: %s", pvName, pkgName)
		allErrs = append(allErrs,
			field.Invalid(fldPath, pvName, "must begin with <spec.packageName> + '.'"))
	}

	return allErrs
}

// validate spec.version is not empty
func ValidatePackageVersionSpecVersion(version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if version == "" {
		allErrs = append(allErrs,
			field.Invalid(fldPath, version, "cannot be empty"))
	}

	return allErrs
}

// validate spec.PackageName isnt empty
func ValidatePackageVersionSpecPackageName(name string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if name == "" {
		allErrs = append(allErrs, field.Required(fldPath, "can not be empty"))
	}

	allErrs = append(allErrs, ValidatePackageName(name, fldPath)...)
	return allErrs
}

func IsFullyQualifiedName(fldPath *field.Path, name string) field.ErrorList {
	var allErrors field.ErrorList
	if len(name) == 0 {
		return append(allErrors, field.Required(fldPath, ""))
	}
	if errs := validation.IsDNS1123Subdomain(name); len(errs) > 0 {
		return append(allErrors, field.Invalid(fldPath, name, strings.Join(errs, ",")))
	}
	if len(strings.Split(name, ".")) < 3 {
		return append(allErrors, field.Invalid(fldPath, name, "should be a domain with at least three segments separated by dots"))
	}
	return allErrors
}