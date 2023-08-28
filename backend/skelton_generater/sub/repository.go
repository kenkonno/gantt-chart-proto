package sub

type Repository struct {
}

func (r *Repository) GetPackage() string {
	return `package repository
`
}

func (r *Repository) GetImports() string {
	return `import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
`
}

// GetConstructor New と struct 定義
func (r *Repository) GetConstructor(structName string) string {
	template :=
		`func New@Upper@Repository() @Lower@Repository {
	return @Lower@Repository{con}
}

type @Lower@Repository struct {
	con *gorm.DB
}`

	return RewriteString(template, structName)

}

// GetDefaultFunctions select, update, delete を追加する
func (r *Repository) GetDefaultFunctions(structName string) string {
	template := `
func (r *@Lower@Repository) FindAll() []db.@Upper@ {
	var @Lower@s []db.@Upper@

	result := r.con.Order("id DESC").Find(&@Lower@s)
	if result.Error != nil {
		panic(result.Error)
	}
	return @Lower@s
}

func (r *@Lower@Repository) Find(id int32) db.@Upper@ {
	var @Lower@ db.@Upper@

	result := r.con.First(&@Lower@, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return @Lower@
}

func (r *@Lower@Repository) Upsert(m db.@Upper@) db.@Upper@ {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *@Lower@Repository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.@Upper@{})
}
`
	return RewriteString(template, structName)
}
