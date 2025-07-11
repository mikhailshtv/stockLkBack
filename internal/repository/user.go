package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/utils/jwtgen"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UsersRepository struct {
	Users    []model.User
	UsersLen int
	db       *sqlx.DB
	redis    *redis.Client
}

func NewUsersRepository(db *sqlx.DB, redis *redis.Client) *UsersRepository {
	return &UsersRepository{db: db, redis: redis}
}

func (ur *UsersRepository) Create(user model.User, ctx context.Context) (*model.User, error) {
	const query = `
		INSERT INTO users.users (
			login, 
			password_hash, 
			first_name, 
			last_name, 
			email
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING *
	`

	if user.PasswordHash == "" {
		return nil, fmt.Errorf("хэш пароля отсутствует")
	}

	err := ur.db.QueryRowxContext(
		ctx,
		query,
		user.Login,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Email,
	).StructScan(&user)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_login_key":
				return nil, fmt.Errorf("логин %s уже существует", user.Login)
			case "users_email_key":
				return nil, fmt.Errorf("email %s уже существует", user.Email)
			default:
				return nil, fmt.Errorf("найден дубликат записи: %w", err)
			}
		}
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	return &user, nil
}

func (ur *UsersRepository) GetAll(ctx context.Context) ([]model.User, error) {
	const query = "SELECT * FROM users.users ORDER BY id"

	var users []model.User
	err := ur.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %w", err)
	}

	return users, nil

}

func (ur *UsersRepository) GetByID(id int, ctx context.Context) (*model.User, error) {
	const query = `
        SELECT 
            id,
            login,
            first_name,
            last_name,
            email,
            role
        FROM users.users
        WHERE id = $1
        LIMIT 1
    `

	var user model.User
	err := ur.db.QueryRowxContext(ctx, query, id).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return &user, nil
}

func (ur *UsersRepository) Delete(id int, ctx context.Context) (*model.User, error) {
	const query = `
        WITH deleted AS (
            DELETE FROM users.users 
            WHERE id = $1
            RETURNING *
        )
        SELECT * FROM deleted
	`

	deletedUser := model.User{}
	err := ur.db.QueryRowxContext(ctx, query, id).StructScan(&deletedUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с ID %d не найден: %w", id, err)
		}
		return nil, fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	return &deletedUser, nil
}

func (ur *UsersRepository) Update(id int, userReq model.UserEditBody, ctx context.Context) (*model.User, error) {
	existingUser, err := ur.GetByID(id, ctx)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	if userReq.FirstName != "" {
		existingUser.FirstName = userReq.FirstName
	}
	if userReq.LastName != "" {
		existingUser.LastName = userReq.LastName
	}
	if userReq.Email != "" {
		existingUser.Email = userReq.Email
	}

	const query = `
        UPDATE users.users SET
            first_name = $1,
            last_name = $2,
            email = $3
        WHERE id = $4
        RETURNING *
    `

	updatedUser := model.User{}
	err = ur.db.QueryRowxContext(
		ctx,
		query,
		existingUser.FirstName,
		existingUser.LastName,
		existingUser.Email,
		id,
	).StructScan(&updatedUser)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				return nil, fmt.Errorf("email %s уже используется другим пользователем", userReq.Email)
			}
		}
		return nil, fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	return &updatedUser, nil
}

func (ur *UsersRepository) Login(userReq model.LoginRequest, ctx context.Context) (*model.TokenSuccess, error) {
	const query = `
        SELECT
						id,
            login,
            password_hash,
            email,
            role
        FROM users.users
        WHERE login = $1 OR email = $1
        LIMIT 1
    `

	var user model.User
	err := ur.db.QueryRowxContext(ctx, query, userReq.Login).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("логин или пароль пользователя недействителен")
		}
		return nil, fmt.Errorf("ошибка запроса к базе: %w", err)
	}
	if user.CheckUserPassword(userReq.Password) {
		token, err := jwtgen.GenerateToken(user.ID, user.Login, user.Role)
		if err != nil {
			return nil, errors.New("ошибка генерации токена")
		}
		return &model.TokenSuccess{
			Message: "аутентификация успешна",
			Token:   token,
		}, nil
	}
	return nil, errors.New("логин или пароль пользователя недействителен")
}

func (ur *UsersRepository) ChangeUserRole(id int, userRoleReq model.UserRoleBody, ctx context.Context) (*model.User, error) {
	if !userRoleReq.Role.Valid() {
		return nil, fmt.Errorf("недопустимая роль: %s", userRoleReq.Role)
	}
	const query = `
		UPDATE users.users SET
				role = $1
		WHERE id = $2
		RETURNING *
	`
	var updatedUser model.User
	err := ur.db.QueryRowxContext(ctx, query, userRoleReq.Role, id).StructScan(&updatedUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка при обновлении роли: %w", err)
	}
	return &updatedUser, nil
}

func (ur *UsersRepository) ChangePassword(
	id int,
	changePassworReq model.UserChangePasswordBody,
	ctx context.Context,
) (*model.Success, error) {
	if changePassworReq.Password != changePassworReq.PasswordConfirm {
		return nil, fmt.Errorf("пароли не совпадают")
	}

	var currentHash string
	const passwordHashQuery = "SELECT password_hash FROM users.users WHERE id = $1"
	err := ur.db.QueryRowContext(ctx, passwordHashQuery, id).Scan(&currentHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	if changePassworReq.OldPassword == "" {
		return nil, fmt.Errorf("необходимо указать текущий пароль")
	}
	var user model.User
	user.PasswordHash = currentHash
	res := user.CheckUserPassword(changePassworReq.OldPassword)
	fmt.Printf("%v", res)
	if !res {
		return nil, fmt.Errorf("неверный текущий пароль")
	}

	err = user.HashPassword(changePassworReq.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}
	const query = `UPDATE users.users SET password_hash = $1 WHERE id = $2`
	result, err := ur.db.ExecContext(ctx, query, user.PasswordHash, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при изменении пароля: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке изменения пароля: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("пользователь не найден")
	}
	return &model.Success{
		Status:  "Success",
		Message: "Пароль успешно изменен",
	}, nil
}

func (or *UsersRepository) WriteLog(result any, operation, status, tableName string) (int64, error) {
	return WriteLog(result, operation, status, tableName, or.redis)
}
