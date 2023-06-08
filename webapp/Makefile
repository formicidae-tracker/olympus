colon := :
$(colon) := :

all:
	npm install
	npm run ng build -- --configuration production --source-map
	npm run ng run olympus$(:)app-shell$(:)production
