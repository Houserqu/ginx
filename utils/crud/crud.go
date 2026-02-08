package crud

import (
	"encoding/json"
	"fmt"
	"ginx/core"
	"ginx/utils"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BuildQueryConditions 根据参数结构体和  标签构建查询条件
func BuildQueryConditions(query *gorm.DB, params interface{}) *gorm.DB {
	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)

	// 如果是指针，获取指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// 遍历结构体的每个字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 获取 form 标签作为字段名
		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		// 获取  标签
		Tag := field.Tag.Get("")

		// 检查字段值是否为零值
		if isZeroValue(value) {
			continue
		}

		// 转换字段名为数据库列名（下划线格式）
		// 优先使用 column tag，如果没有则使用 form tag 转换
		columnName := field.Tag.Get("column")
		if columnName == "" {
			columnName = toSnakeCase(formTag)
		}

		// 根据  标签类型构建查询条件
		switch Tag {
		case "search":
			// 多字段 OR LIKE 组合查询
			searchFields := field.Tag.Get("search_fields")
			if searchFields == "" {
				// 如果没有指定 search_fields，使用当前字段名
				searchFields = formTag
			}

			fields := strings.Split(searchFields, ",")
			if len(fields) > 0 {
				var conditions []string
				var args []interface{}
				searchValue := "%" + fmt.Sprint(value.Interface()) + "%"

				for _, f := range fields {
					f = strings.TrimSpace(f)
					columnName := toSnakeCase(f)
					conditions = append(conditions, fmt.Sprintf("`%s` LIKE ?", columnName))
					args = append(args, searchValue)
				}

				// 使用 OR 连接所有条件
				orCondition := strings.Join(conditions, " OR ")
				query = query.Where(orCondition, args...)
			}
		case "like":
			// 模糊查询
			query = query.Where(fmt.Sprintf("%s LIKE ?", columnName), "%"+fmt.Sprint(value.Interface())+"%")
		case "left_like":
			// 左模糊查询
			query = query.Where(fmt.Sprintf("%s LIKE ?", columnName), "%"+fmt.Sprint(value.Interface()))
		case "right_like":
			// 右模糊查询
			query = query.Where(fmt.Sprintf("%s LIKE ?", columnName), fmt.Sprint(value.Interface())+"%")
		case "in":
			// IN 查询
			query = query.Where(fmt.Sprintf("%s IN ?", columnName), value.Interface())
		case "not_in":
			// NOT IN 查询
			query = query.Where(fmt.Sprintf("%s NOT IN ?", columnName), value.Interface())
		case "gt":
			// 大于
			query = query.Where(fmt.Sprintf("%s > ?", columnName), value.Interface())
		case "gte":
			// 大于等于
			query = query.Where(fmt.Sprintf("%s >= ?", columnName), value.Interface())
		case "lt":
			// 小于
			query = query.Where(fmt.Sprintf("%s < ?", columnName), value.Interface())
		case "lte":
			// 小于等于
			query = query.Where(fmt.Sprintf("%s <= ?", columnName), value.Interface())
		case "ne":
			// 不等于
			query = query.Where(fmt.Sprintf("%s != ?", columnName), value.Interface())
		case "between":
			// BETWEEN 查询（需要传入数组）
			if value.Kind() == reflect.Slice && value.Len() == 2 {
				query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", columnName), value.Index(0).Interface(), value.Index(1).Interface())
			}
		case "in_set":
			// FIND_IN_SET 查询
			query = query.Where(fmt.Sprintf("FIND_IN_SET(?, %s)", columnName), value.Interface())
		default:
			// 默认为等于查询
			query = query.Where(fmt.Sprintf("%s = ?", columnName), value.Interface())
		}
	}

	return query
}

// isZeroValue 检查值是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// toSnakeCase 将驼峰命名转换为下划线命名
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

type Options struct {
	Select   []string // select 字段
	PreLoads []string
}

type Option func(*Options)

func WithSelect(fields ...string) Option {
	return func(o *Options) {
		o.Select = fields
	}
}

func WithPreLoads(preloads ...string) Option {
	return func(o *Options) {
		o.PreLoads = preloads
	}
}

func List[M any](c *gin.Context, params any, opts ...Option) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	query := utils.MysqlC(c).Model(new(M))
	if err := c.ShouldBindQuery(params); err != nil {
		core.Error(c, err.Error())
		return
	}

	v := reflect.ValueOf(new(M)).Elem()

	// 应用查询条件
	query = BuildQueryConditions(query, params)

	// 判断是否存在 TeamId 字段，若存在则添加团队隔离查询条件
	teamIdField := v.FieldByName("TeamId")
	if teamIdField.IsValid() {
		query = query.Where("team_id = ?", c.GetInt64("TeamId"))
	}

	// 查总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		core.Error(c, err.Error())
		return
	}

	// 选择字段
	if len(options.Select) > 0 {
		query = query.Select(options.Select)
	}

	list := []M{}
	for _, preLoad := range options.PreLoads {
		query = query.Preload(preLoad)
	}

	if err := query.Scopes(utils.Paginate(c)).Order("id DESC").Find(&list).Error; err != nil {
		core.Error(c, err.Error())
		return
	}

	core.Success(c, gin.H{
		"total": total,
		"list":  list,
		"page":  utils.GetPage(c),
		"size":  utils.GetPageSize(c),
	})
}

func DetailByID[M any](c *gin.Context, opts ...Option) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	query := utils.MysqlC(c).Model(new(M))

	id := c.Query("id")
	if id == "" {
		core.Error(c, "id不能为空")
		return
	}
	query = query.Where("id = ?", id)

	v := reflect.ValueOf(new(M)).Elem()

	// 判断是否存在 TeamId 字段，若存在则添加团队隔离查询条件
	teamIdField := v.FieldByName("TeamId")
	if teamIdField.IsValid() {
		query = query.Where("team_id = ?", c.GetInt64("TeamId"))
	}

	for _, preLoad := range options.PreLoads {
		query = query.Preload(preLoad)
	}

	var detail M
	if err := query.Take(&detail).Error; err != nil {
		core.Error(c, err.Error())
		return
	}

	core.Success(c, detail)
	return
}

func DeleteByID[M any](c *gin.Context) {
	var params struct {
		Id int64 `json:"id,string" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		core.Error(c, err.Error())
		return
	}
	query := utils.MysqlC(c).Model(new(M))

	v := reflect.ValueOf(new(M)).Elem()

	// 判断是否存在 TeamId 字段，若存在则添加团队隔离查询条件
	teamIdField := v.FieldByName("TeamId")
	if teamIdField.IsValid() {
		query = query.Where("team_id = ?", c.GetInt64("TeamId"))
	}

	if err := query.Where("id = ?", params.Id).Delete(nil).Error; err != nil {
		core.Error(c, err.Error())
		return
	}

	core.Success(c, "删除成功")
}

func Update[M any](c *gin.Context, params any) {
	// 先读取原始请求体，以便判断哪些字段在 JSON 中被显式传入（包含零值）
	body, err := c.GetRawData()
	if err != nil {
		core.Error(c, err.Error())
		return
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(body, &rawMap); err != nil {
		core.Error(c, err.Error())
		return
	}

	// 将请求体反序列化到 params 结构体，保持类型信息
	if err := json.Unmarshal(body, params); err != nil {
		core.Error(c, err.Error())
		return
	}

	v := reflect.ValueOf(params).Elem()
	t := v.Type()

	idField := v.FieldByName("Id")
	if !idField.IsValid() {
		core.Error(c, "参数结构体必须包含 Id 字段")
		return
	}
	id := idField.Int()

	updates := make(map[string]interface{})

	// 遍历结构体字段，依据 JSON key 是否存在来决定是否包含该字段（允许零值）
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Id" {
			continue
		}

		// 优先取 json tag
		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]
		if jsonKey == "" || jsonKey == "-" {
			// fallback 到 form tag
			jsonKey = field.Tag.Get("form")
		}
		if jsonKey == "" {
			// 最后 fallback 到首字母小写的字段名
			jsonKey = strings.ToLower(field.Name[:1]) + field.Name[1:]
		}

		// 如果该字段在请求 JSON 中显式出现，则包含到 updates（哪怕是零值）
		if _, ok := rawMap[jsonKey]; ok {
			columnName := toSnakeCase(field.Name)
			updates[columnName] = v.Field(i).Interface()
		}
	}

	if len(updates) == 0 {
		core.Error(c, "没有可更新的字段")
		return
	}

	query := utils.MysqlC(c).Model(new(M)).Where("id = ?", id)

	// 仅当模型 M 存在 TeamId 字段时添加团队隔离查询条件
	modelV := reflect.ValueOf(new(M)).Elem()
	if modelV.FieldByName("TeamId").IsValid() {
		query = query.Where("team_id = ?", c.GetInt64("TeamId"))
	}

	if err := query.Updates(updates).Error; err != nil {
		core.Error(c, err.Error())
		return
	}

	core.Success(c, updates)
}

func Create[M any](c *gin.Context, params any) {
	if err := c.ShouldBindJSON(params); err != nil {
		core.Error(c, err.Error())
		return
	}

	var created M

	createdValue := reflect.ValueOf(&created).Elem()
	paramsValue := reflect.ValueOf(params).Elem()
	paramsType := reflect.TypeOf(params).Elem()

	// 将 params 的字段值赋值到 created 对应字段
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		paramFieldValue := paramsValue.Field(i)

		createdFieldValue := createdValue.FieldByName(field.Name)
		if createdFieldValue.IsValid() && createdFieldValue.CanSet() {
			createdFieldValue.Set(paramFieldValue)
		}
	}

	// 生成雪花 ID
	idField := createdValue.FieldByName("Id")
	if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.Int64 {
		idField.SetInt(utils.ID())
	}

	// 如果模型中存在 TeamId 字段，从 ctx 中获取
	teamIdField := createdValue.FieldByName("TeamId")
	if teamIdField.IsValid() && teamIdField.CanSet() && teamIdField.Kind() == reflect.Int64 {
		if teamId, exists := c.Get("TeamId"); exists {
			if id, ok := teamId.(int64); ok {
				teamIdField.SetInt(id)
			}
		}
	}

	// 如果模型中存在 CreatorId 字段，从 ctx 中获取
	creatorIdField := createdValue.FieldByName("CreatorId")
	if creatorIdField.IsValid() && creatorIdField.CanSet() && creatorIdField.Kind() == reflect.Int64 {
		if creatorId, exists := c.Get("UserId"); exists {
			if id, ok := creatorId.(int64); ok {
				creatorIdField.SetInt(id)
			}
		}
	}

	if err := utils.MysqlC(c).Create(&created).Error; err != nil {
		core.Error(c, err.Error())
		return
	}

	core.Success(c, created)
}
