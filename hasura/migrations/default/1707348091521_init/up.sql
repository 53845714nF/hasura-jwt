SET check_function_bodies = false;
CREATE FUNCTION public.create_assigned_user_role() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN 
    INSERT INTO assigned_user_roles (user_id, user_role_name)
    VALUES (NEW.id, 'user');
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.set_current_timestamp_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$;
CREATE TABLE public.assigned_user_roles (
    user_id uuid NOT NULL,
    user_role_name text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public."user" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.user_roles (
    name text NOT NULL
);
ALTER TABLE ONLY public.assigned_user_roles
    ADD CONSTRAINT assigned_user_roles_pkey PRIMARY KEY (user_id, user_role_name);
ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_key UNIQUE (email);
ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (name);
CREATE TRIGGER on_create_user AFTER INSERT ON public."user" FOR EACH ROW EXECUTE FUNCTION public.create_assigned_user_role();
CREATE TRIGGER set_public_assigned_user_roles_updated_at BEFORE UPDATE ON public.assigned_user_roles FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_assigned_user_roles_updated_at ON public.assigned_user_roles IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER set_public_user_updated_at BEFORE UPDATE ON public."user" FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER set_public_user_updated_at ON public."user" IS 'trigger to set value of column "updated_at" to current timestamp on row update';
ALTER TABLE ONLY public.assigned_user_roles
    ADD CONSTRAINT "assigned_user_roles_User_role_name_fkey" FOREIGN KEY (user_role_name) REFERENCES public.user_roles(name) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.assigned_user_roles
    ADD CONSTRAINT assigned_user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
