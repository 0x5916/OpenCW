import * as m from '$lib/paraglide/messages';

function isErrorCode(value: string): boolean {
	return /^[A-Z0-9_]+$/.test(value);
}

function readCodeFromBody(body: unknown): string | null {
	if (!body || typeof body !== 'object') return null;

	const code = (body as Record<string, unknown>).code;
	if (typeof code === 'string' && code.trim() !== '') return code;

	const error = (body as Record<string, unknown>).error;
	if (typeof error === 'string' && isErrorCode(error.trim())) return error;

	const message = (body as Record<string, unknown>).message;
	if (typeof message === 'string' && isErrorCode(message.trim())) return message;

	return null;
}

function readErrorCode(error: unknown): string | null {
	if (!error || typeof error !== 'object') return null;

	const code = (error as { code?: unknown }).code;
	if (typeof code === 'string' && code.trim() !== '') return code;

	const bodyCode = readCodeFromBody((error as { body?: unknown }).body);
	if (bodyCode) return bodyCode;

	if (error instanceof Error && isErrorCode(error.message.trim())) {
		return error.message;
	}

	return null;
}

export function localizeApiError(error: unknown, fallback: () => string): string {
	const code = readErrorCode(error);

	switch (code) {
		case 'INVALID_REQUEST_BODY':
			return m.api_error_invalid_request_body();
		case 'INTERNAL_SERVER_ERROR':
			return m.api_error_internal_server_error();
		case 'DATABASE_FAILURE':
			return m.api_error_database_failure();
		case 'INVALID_CREDENTIALS':
			return m.api_error_invalid_credentials();
		case 'CONFLICT':
			return m.api_error_conflict();
		case 'INVALID_TOKEN':
			return m.api_error_invalid_token();
		case 'EXPIRED_TOKEN':
			return m.api_error_expired_token();
		case 'AUTH_HEADER_REQUIRED':
			return m.api_error_auth_header_required();
		case 'INVALID_AUTH_HEADER_FORMAT':
			return m.api_error_invalid_auth_header_format();
		case 'USER_NOT_FOUND':
			return m.api_error_user_not_found();
		case 'SETTINGS_FETCH_FAILED':
			return m.api_error_settings_fetch_failed();
		case 'SETTINGS_UPDATE_FAILED':
			return m.api_error_settings_update_failed();
		case 'PROGRESS_QUERY_FAILED':
			return m.api_error_progress_query_failed();
		case 'PROGRESS_CREATE_FAILED':
			return m.api_error_progress_create_failed();
		case 'PASSWORD_HASH_FAILED':
			return m.api_error_password_hash_failed();
		case 'TOKEN_ISSUE_FAILED':
			return m.api_error_token_issue_failed();
		case 'EMAIL_ALREADY_IN_USE':
			return m.api_error_email_already_in_use();
		case 'EMAIL_VERIFIED_BY_ANOTHER_ACCOUNT':
			return m.api_error_email_verified_by_another_account();
		case 'USERNAME_ALREADY_IN_USE':
			return m.api_error_username_already_in_use();
		case 'USERNAME_EMAIL_ALREADY_IN_USE':
			return m.api_error_username_email_already_in_use();
		case 'EMAIL_UNCHANGED':
			return m.api_error_email_unchanged();
		case 'CALL_SIGN_ALREADY_IN_USE':
			return m.api_error_call_sign_already_in_use();
		case 'EMAIL_ALREADY_VERIFIED':
			return m.api_error_email_already_verified();
		case 'VERIFICATION_RATE_LIMITED':
			return m.api_error_verification_rate_limited();
		case 'VERIFICATION_SEND_FAILED':
			return m.api_error_verification_send_failed();
		case 'VERIFICATION_CODE_INVALID':
			return m.api_error_verification_code_invalid();
		case 'VERIFICATION_CODE_EXPIRED':
			return m.api_error_verification_code_expired();
		case 'VALIDATION_FAILED':
			return m.api_error_validation_failed();
		case 'LOGIN_FAILED':
			return m.api_error_login_failed();
		case 'REGISTER_FAILED':
			return m.api_error_register_failed();
		default:
			return fallback();
	}
}
