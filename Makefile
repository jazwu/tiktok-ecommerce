.PHONY: gen-frontend
gen-frontend:
	cwgo server --type HTTP --idl ../../idl/frontend/home.proto --service frontend -module github.com/jazwu/tiktok-ecommerce/app/f
rontend -I ../../idl